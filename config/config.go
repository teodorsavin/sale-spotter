package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"teodorsavin/ah-bonus/model"
)

func SetupDatabase() (*gorm.DB, error) {
	// Validate environment variables
	host, name, user, password := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")
	if host == "" || name == "" || user == "" || password == "" {
		return nil, fmt.Errorf("missing required environment variables for database connection")
	}

	// Build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, name)

	// Connect to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Set collation for all tables
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	// Run migrations
	err = db.AutoMigrate(&model.Product{}, &model.Image{}, &model.DiscountLabel{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
