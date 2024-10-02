package main

import (
	"github.com/gin-gonic/gin"

	"teodorsavin/ah-bonus/controller"
)

func main() {
	router := gin.Default()

	// Routes requiring authentication
	privateRouter := router.Group("/api")
	privateRouter.Use(controller.Authenticate())
	{
		privateRouter.GET("/products", controller.GetProducts)
		privateRouter.GET("/brands", controller.GetBrands)
	}

	// Login endpoint
	router.POST("/login", controller.Login)

	// Test endpoint to save products from test.json into the database
	router.GET("/test/saveproducts", controller.SaveProductsDebug)

	router.Run("0.0.0.0:8080")
}
