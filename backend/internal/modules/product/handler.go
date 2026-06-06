package product

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"medcart-backend/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const cacheKey = "products_list"

type Handler struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewHandler(db *sql.DB, rdb *redis.Client) *Handler {
	return &Handler{DB: db, Redis: rdb}
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"imageUrl"`
}

type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"imageUrl"`
}

func (h *Handler) clearCache() {
	if h.Redis != nil {
		h.Redis.Del(cache.Ctx, cacheKey)
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	if h.Redis != nil {
		cached, err := h.Redis.Get(cache.Ctx, cacheKey).Result()
		if err == nil {
			var products []Product
			if json.Unmarshal([]byte(cached), &products) == nil {
				c.JSON(http.StatusOK, gin.H{"products": products, "source": "cache"})
				return
			}
		}
	}

	rows, err := h.DB.Query("SELECT id, name, description, category, price, stock, image_url FROM products ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch products"})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		var description, imageURL sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &description, &p.Category, &p.Price, &p.Stock, &imageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read products"})
			return
		}
		p.Description = description.String
		p.ImageURL = imageURL.String
		products = append(products, p)
	}

	if h.Redis != nil {
		data, err := json.Marshal(products)
		if err == nil {
			h.Redis.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
		}
	}

	c.JSON(http.StatusOK, gin.H{"products": products, "source": "database"})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var p Product
	var description, imageURL sql.NullString
	err = h.DB.QueryRow(
		"SELECT id, name, description, category, price, stock, image_url FROM products WHERE id = $1", id,
	).Scan(&p.ID, &p.Name, &description, &p.Category, &p.Price, &p.Stock, &imageURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	p.Description = description.String
	p.ImageURL = imageURL.String

	c.JSON(http.StatusOK, gin.H{"product": p})
}

func (h *Handler) Create(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Name == "" || req.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and category are required"})
		return
	}

	var id int
	err := h.DB.QueryRow(
		"INSERT INTO products (name, description, category, price, stock, image_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Name, req.Description, req.Category, req.Price, req.Stock, req.ImageURL,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create product"})
		return
	}

	h.clearCache()
	c.JSON(http.StatusCreated, gin.H{"message": "product created", "id": id})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE products SET name = $1, description = $2, category = $3, price = $4, stock = $5, image_url = $6 WHERE id = $7",
		req.Name, req.Description, req.Category, req.Price, req.Stock, req.ImageURL, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product"})
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	h.clearCache()
	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	h.clearCache()
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
