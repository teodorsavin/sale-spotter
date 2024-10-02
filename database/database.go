package database

import (
	"gorm.io/gorm"
	"teodorsavin/ah-bonus/model"
)

// AllBrands retrieves unique brands from the Product table
func AllBrands(db *gorm.DB) ([]string, error) {
	var products []model.Product
	var brands []string

	// Query all products
	if err := db.Select("DISTINCT brand").Find(&products).Error; err != nil {
		return nil, err
	}

	// Extract unique brands
	for _, product := range products {
		if !contains(brands, product.Brand) {
			brands = append(brands, product.Brand)
		}
	}

	return brands, nil
}

// Helper function to check if a slice contains a value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// InsertProductsBulk inserts products in bulk into the database
func InsertProductsBulk(db *gorm.DB, products []model.Product) error {
	return db.Create(&products).Error
}
