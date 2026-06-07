package product

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"medcart-backend/internal/cache"
	"medcart-backend/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const cacheKey = "products_list"

type Handler struct {
	DB    *mongo.Database
	Redis *redis.Client
}

func NewHandler(db *mongo.Database, rdb *redis.Client) *Handler {
	return &Handler{DB: db, Redis: rdb}
}

type Product struct {
	ID          int     `bson:"id" json:"id"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	Category    string  `bson:"category" json:"category"`
	Price       float64 `bson:"price" json:"price"`
	Stock       int     `bson:"stock" json:"stock"`
	ImageURL    string  `bson:"image_url" json:"imageUrl"`
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

	cursor, err := h.DB.Collection("products").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch products"})
		return
	}
	defer cursor.Close(context.Background())

	products := []Product{}
	if err := cursor.All(context.Background(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read products"})
		return
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
	err = h.DB.Collection("products").FindOne(context.Background(), bson.M{"id": id}).Decode(&p)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

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

	id, err := database.NextID(h.DB, "products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create product id"})
		return
	}

	product := Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}
	_, err = h.DB.Collection("products").InsertOne(context.Background(), product)
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

	update := bson.M{"$set": bson.M{
		"name":        req.Name,
		"description": req.Description,
		"category":    req.Category,
		"price":       req.Price,
		"stock":       req.Stock,
		"image_url":   req.ImageURL,
	}}
	result, err := h.DB.Collection("products").UpdateOne(context.Background(), bson.M{"id": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product"})
		return
	}
	if result.MatchedCount == 0 {
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

	result, err := h.DB.Collection("products").DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	h.clearCache()
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
