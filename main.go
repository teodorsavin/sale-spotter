package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"teodorsavin/ah-bonus/config"
	"teodorsavin/ah-bonus/controller"
)

func main() {
	db, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}

	router := gin.Default()

	// Routes requiring authentication
	privateRouter := router.Group("/api")
	privateRouter.Use(controller.Authenticate())
	{
		privateRouter.GET("/products", controller.GetProducts(db))
		privateRouter.GET("/brands", controller.GetBrands(db))
	}

	// Login endpoint
	router.POST("/login", controller.Login)

	// Test endpoint to save products from test.json into the database
	router.GET("/test/saveproducts", controller.SaveProductsDebug(db))

	router.Run("0.0.0.0:8080")
}
