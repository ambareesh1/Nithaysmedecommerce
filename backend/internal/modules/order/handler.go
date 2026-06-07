package order

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"medcart-backend/internal/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	DB *mongo.Database
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{DB: db}
}

type OrderItemRequest struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type CreateRequest struct {
	Items []OrderItemRequest `json:"items"`
}

type OrderItem struct {
	ProductID int     `bson:"product_id" json:"productId"`
	Name      string  `bson:"name" json:"name"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	Price     float64 `bson:"price" json:"price"`
	Subtotal  float64 `bson:"subtotal" json:"subtotal"`
}

type Order struct {
	ID            int         `bson:"id" json:"id"`
	UserID        int         `bson:"user_id" json:"-"`
	OrderNumber   string      `bson:"order_number" json:"orderNumber"`
	TotalAmount   float64     `bson:"total_amount" json:"totalAmount"`
	TotalQuantity int         `bson:"total_quantity" json:"totalQuantity"`
	Status        string      `bson:"status" json:"status"`
	CreatedAt     time.Time   `bson:"created_at" json:"createdAt"`
	Items         []OrderItem `bson:"items" json:"items"`
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetInt("userID")
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order must have items"})
		return
	}

	totalAmount := 0.0
	totalQuantity := 0
	items := []OrderItem{}

	for _, item := range req.Items {
		var product struct {
			Name  string  `bson:"name"`
			Price float64 `bson:"price"`
		}
		err := h.DB.Collection("products").FindOne(context.Background(), bson.M{"id": item.ProductID}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}
		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal
		totalQuantity += item.Quantity
		items = append(items, OrderItem{
			ProductID: item.ProductID,
			Name:      product.Name,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		})
	}

	orderID, err := database.NextID(h.DB, "orders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create order id"})
		return
	}

	orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())

	order := Order{
		ID:            orderID,
		UserID:        userID,
		OrderNumber:   orderNumber,
		TotalAmount:   totalAmount,
		TotalQuantity: totalQuantity,
		Status:        "pending",
		CreatedAt:     time.Now(),
		Items:         items,
	}
	_, err = h.DB.Collection("orders").InsertOne(context.Background(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create order"})
		return
	}

	h.DB.Collection("carts").DeleteMany(context.Background(), bson.M{"user_id": userID})

	c.JSON(http.StatusCreated, gin.H{
		"message":     "order placed",
		"id":          orderID,
		"orderNumber": orderNumber,
	})
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	userID := c.GetInt("userID")

	opts := options.Find().SetSort(bson.M{"id": -1})
	cursor, err := h.DB.Collection("orders").Find(context.Background(), bson.M{"user_id": userID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch orders"})
		return
	}
	defer cursor.Close(context.Background())

	orders := []Order{}
	if err := cursor.All(context.Background(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var o Order
	err = h.DB.Collection("orders").FindOne(context.Background(), bson.M{"id": id}).Decode(&o)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": o})
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.DB.Collection("orders").UpdateOne(
		context.Background(),
		bson.M{"id": id},
		bson.M{"$set": bson.M{"status": req.Status}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
