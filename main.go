package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"teodorsavin/ah-bonus/config"
	"teodorsavin/ah-bonus/controller"
)

func main() {
	app := fx.New(
		// Provide the dependencies to be injected
		fx.Provide(
			config.SetupDatabase,
			NewRouter,
		),
		// Invoke will run the main application logic
		fx.Invoke(registerRoutes),
	)

	app.Run()
}

func NewRouter() *gin.Engine {
	return gin.Default()
}

func registerRoutes(lc fx.Lifecycle, router *gin.Engine, db *gorm.DB) {
	// Routes requiring authentication
	privateRouter := router.Group("/api")
	privateRouter.Use(controller.Authenticate())
	{
		privateRouter.GET("/products", controller.GetProducts(db))
		privateRouter.GET("/brands", controller.GetBrands(db))
	}

	router.POST("/login", controller.Login)

	// Test endpoint to save products from test.json into the database
	router.GET("/test/saveproducts", controller.SaveProductsDebug(db))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := router.Run("0.0.0.0:8080"); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down the server...")
			return nil
		},
	})
}
