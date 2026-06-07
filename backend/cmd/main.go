package main

import (
	"log"
	"net/http"

	"medcart-backend/internal/cache"
	"medcart-backend/internal/config"
	"medcart-backend/internal/database"
	"medcart-backend/internal/middleware"
	"medcart-backend/internal/modules/auth"
	"medcart-backend/internal/modules/cart"
	"medcart-backend/internal/modules/order"
	"medcart-backend/internal/modules/product"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}

	rdb := cache.Connect(cfg.RedisAddr)

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := auth.NewHandler(db, cfg.JWTSecret)
	productHandler := product.NewHandler(db, rdb)
	cartHandler := cart.NewHandler(db)
	orderHandler := order.NewHandler(db)

	api := router.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		api.GET("/products", productHandler.GetAll)
		api.GET("/products/:id", productHandler.GetByID)
		api.POST("/products", productHandler.Create)
		api.PUT("/products/:id", productHandler.Update)
		api.DELETE("/products/:id", productHandler.Delete)

		protected := api.Group("")
		protected.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			protected.GET("/cart", cartHandler.GetCart)
			protected.POST("/cart", cartHandler.AddToCart)
			protected.PUT("/cart/:productId", cartHandler.UpdateItem)
			protected.DELETE("/cart/:productId", cartHandler.RemoveItem)

			protected.POST("/orders", orderHandler.Create)
			protected.GET("/orders/my", orderHandler.GetMyOrders)
			protected.GET("/orders/:id", orderHandler.GetByID)
			protected.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
		}
	}

	log.Println("server running on port " + cfg.Port)
	router.Run(":" + cfg.Port)
}
