package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"teodorsavin/ah-bonus/database"
	"teodorsavin/ah-bonus/model"
)

func SaveProductsDebug(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products := []model.Product{}
		err := database.InsertProductsBulk(db, products)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Products saved successfully!"})
	}
}
