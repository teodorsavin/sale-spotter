package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"teodorsavin/ah-bonus/database"
)

// GetBrands responds with the list of all brands as JSON.
func GetBrands(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		brands, err := database.AllBrands(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, brands)
	}
}
