package controller

import (
	"gorm.io/gorm"
	"teodorsavin/ah-bonus/database"
	"teodorsavin/ah-bonus/model"
	"teodorsavin/ah-bonus/service"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	baseUrl = "https://api.ah.nl"
	timeout = time.Second * 30
)

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []model.Product
		if err := db.Preload("Images").Preload("DiscountLabels").Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(products) == 0 {
			products = RequestProducts(db)
		}

		c.JSON(200, products)
	}
}

func RequestProducts(db *gorm.DB) []model.Product {
	apiClient := service.NewAPIClient(baseUrl, timeout)
	token := apiClient.Login()
	dataFromAPI := apiClient.GetProducts(token, 0)

	if len(dataFromAPI.Products) > 0 {
		err := database.InsertProductsBulk(db, dataFromAPI.Products)
		if err != nil {
			return []model.Product{} // return empty slice if error occurs
		}
	}
	return dataFromAPI.Products
}
