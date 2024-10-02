package main

import (
	"fmt"
	"log"
	"os"
	"teodorsavin/ah-bonus/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"teodorsavin/ah-bonus/controller"
)

func setupDatabase() (*gorm.DB, error) {
	// Get environment variables
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, name)

	// Connect to the database
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to MySQL database!")
	}
	// Set collation for all tables
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	// Run migrations
	db.AutoMigrate(&model.Product{}, &model.Image{}, &model.DiscountLabel{})

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
