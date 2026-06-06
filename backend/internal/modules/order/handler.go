package order

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
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
	ProductID int     `json:"productId"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type Order struct {
	ID            int         `json:"id"`
	OrderNumber   string      `json:"orderNumber"`
	TotalAmount   float64     `json:"totalAmount"`
	TotalQuantity int         `json:"totalQuantity"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"createdAt"`
	Items         []OrderItem `json:"items"`
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
		var name string
		var price float64
		err := h.DB.QueryRow("SELECT name, price FROM products WHERE id = $1", item.ProductID).Scan(&name, &price)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}
		subtotal := price * float64(item.Quantity)
		totalAmount += subtotal
		totalQuantity += item.Quantity
		items = append(items, OrderItem{
			ProductID: item.ProductID,
			Name:      name,
			Quantity:  item.Quantity,
			Price:     price,
			Subtotal:  subtotal,
		})
	}

	orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())

	var orderID int
	err := h.DB.QueryRow(
		"INSERT INTO orders (user_id, order_number, total_amount, total_quantity, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userID, orderNumber, totalAmount, totalQuantity, "pending",
	).Scan(&orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create order"})
		return
	}

	for _, item := range items {
		_, err := h.DB.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price, subtotal) VALUES ($1, $2, $3, $4, $5)",
			orderID, item.ProductID, item.Quantity, item.Price, item.Subtotal,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save order items"})
			return
		}
	}

	h.DB.Exec("DELETE FROM carts WHERE user_id = $1", userID)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "order placed",
		"id":          orderID,
		"orderNumber": orderNumber,
	})
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	userID := c.GetInt("userID")
	rows, err := h.DB.Query(
		"SELECT id, order_number, total_amount, total_quantity, status, created_at FROM orders WHERE user_id = $1 ORDER BY id DESC",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch orders"})
		return
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.OrderNumber, &o.TotalAmount, &o.TotalQuantity, &o.Status, &o.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read orders"})
			return
		}
		o.Items = h.loadItems(o.ID)
		orders = append(orders, o)
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
	err = h.DB.QueryRow(
		"SELECT id, order_number, total_amount, total_quantity, status, created_at FROM orders WHERE id = $1", id,
	).Scan(&o.ID, &o.OrderNumber, &o.TotalAmount, &o.TotalQuantity, &o.Status, &o.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	o.Items = h.loadItems(o.ID)

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

	result, err := h.DB.Exec("UPDATE orders SET status = $1 WHERE id = $2", req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *Handler) loadItems(orderID int) []OrderItem {
	items := []OrderItem{}
	rows, err := h.DB.Query(
		`SELECT oi.product_id, p.name, oi.quantity, oi.price, oi.subtotal
		 FROM order_items oi JOIN products p ON oi.product_id = p.id
		 WHERE oi.order_id = $1`, orderID,
	)
	if err != nil {
		return items
	}
	defer rows.Close()
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Quantity, &item.Price, &item.Subtotal); err == nil {
			items = append(items, item)
		}
	}
	return items
}
