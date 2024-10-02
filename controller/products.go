package controller

import (
	"net/http"
	"teodorsavin/ah-bonus/service"
	"time"

	"github.com/gin-gonic/gin"

	"teodorsavin/ah-bonus/database"
)

const (
	baseUrl = "https://api.ah.nl"
	timeout = time.Second * 30
)

// GetProducts responds with the list of all products as JSON.
func GetProducts(c *gin.Context) {
	data := database.GetAllProducts()

	// If we don't find products in out database, we call the API and save products in our database
	if len(data.Products) == 0 {
		apiClient := service.NewAPIClient(baseUrl, timeout)
		token := apiClient.Login()
		dataFromAPI := apiClient.GetProducts(token, 0)

		if len(dataFromAPI.Products) > 0 {
			err := database.InsertProductsBulk(dataFromAPI.Products)
			if err != nil {
				return
			}
			data = dataFromAPI
		}
	}

	c.IndentedJSON(http.StatusOK, data)
}
