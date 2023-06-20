package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"teodorsavin/ah-bonus/service"
	"time"

	controller "teodorsavin/ah-bonus/controller"
	model "teodorsavin/ah-bonus/model"
)

const (
	baseUrl = "https://api.ah.nl"
	timeout = time.Second * 30
)

// getProducts responds with the list of all products as JSON.
func getProducts(c *gin.Context) {
	data := controller.GetAllProducts()

	// If we don't find products in out database, we call the API and save products in our database
	if len(data.Products) == 0 {
		apiClient := service.NewAPIClient(baseUrl, timeout)
		token := apiClient.Login()
		dataFromAPI := apiClient.GetProducts(token, 0)

		if len(dataFromAPI.Products) > 0 {
			controller.InsertProductsBulk(dataFromAPI.Products)
			data = dataFromAPI
		}
	}

	c.IndentedJSON(http.StatusOK, data)
}

// getBrands responds with the list of all brands as JSON.
func getBrands(c *gin.Context) {
	data := controller.AllBrands()

	c.IndentedJSON(http.StatusOK, data)
}

func saveProductsDebug(c *gin.Context) {
	filename := "../usr/src/app/database-sample/test.json"

	file, _ := os.ReadFile(filename)
	var data model.BonusProducts
	err := json.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("impossible to read/unmarshall test.json: %s", err)
	}
	controller.InsertProductsBulk(data.Products)

	c.IndentedJSON(http.StatusOK, data)
}

func main() {
	router := gin.Default()

	router.GET("/products", getProducts)
	router.GET("/brands", getBrands)

	// to save all products from test.json into the database
	router.GET("/test/saveproducts", saveProductsDebug)

	router.Run("0.0.0.0:8080")
}
