package cart

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

type CartItem struct {
	ID        int     `json:"id"`
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
	rows, err := h.DB.Query(
		`SELECT c.id, c.product_id, p.name, p.price, c.quantity
		 FROM carts c JOIN products p ON c.product_id = p.id
		 WHERE c.user_id = $1 ORDER BY c.id`, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch cart"})
		return
	}
	defer rows.Close()

	items := []CartItem{}
	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Price, &item.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read cart"})
			return
		}
		items = append(items, item)
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

	_, err := h.DB.Exec(
		"INSERT INTO carts (user_id, product_id, quantity) VALUES ($1, $2, $3)",
		userID, req.ProductID, req.Quantity,
	)
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

	_, err = h.DB.Exec(
		"UPDATE carts SET quantity = $1 WHERE user_id = $2 AND product_id = $3",
		req.Quantity, userID, productID,
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

	_, err = h.DB.Exec(
		"DELETE FROM carts WHERE user_id = $1 AND product_id = $2",
		userID, productID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
