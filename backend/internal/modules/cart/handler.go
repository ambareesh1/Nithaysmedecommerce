package cart

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB *mongo.Database
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{DB: db}
}

type CartItem struct {
	ProductID int     `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type AddRequest struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type UpdateRequest struct {
	Quantity int `json:"quantity"`
}

func (h *Handler) GetCart(c *gin.Context) {
	userID := c.GetInt("userID")

	cursor, err := h.DB.Collection("carts").Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch cart"})
		return
	}
	defer cursor.Close(context.Background())

	items := []CartItem{}
	for cursor.Next(context.Background()) {
		var doc struct {
			ProductID int `bson:"product_id"`
			Quantity  int `bson:"quantity"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		var product struct {
			Name  string  `bson:"name"`
			Price float64 `bson:"price"`
		}
		h.DB.Collection("products").FindOne(context.Background(), bson.M{"id": doc.ProductID}).Decode(&product)

		items = append(items, CartItem{
			ProductID: doc.ProductID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  doc.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) AddToCart(c *gin.Context) {
	userID := c.GetInt("userID")
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	_, err := h.DB.Collection("carts").InsertOne(context.Background(), bson.M{
		"user_id":    userID,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "item added to cart"})
}

func (h *Handler) UpdateItem(c *gin.Context) {
	userID := c.GetInt("userID")
	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err = h.DB.Collection("carts").UpdateOne(
		context.Background(),
		bson.M{"user_id": userID, "product_id": productID},
		bson.M{"$set": bson.M{"quantity": req.Quantity}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart updated"})
}

func (h *Handler) RemoveItem(c *gin.Context) {
	userID := c.GetInt("userID")
	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	_, err = h.DB.Collection("carts").DeleteOne(
		context.Background(),
		bson.M{"user_id": userID, "product_id": productID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
