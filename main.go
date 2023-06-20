package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
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

	//if data.Products == nil {
	//	apiClient := service.NewAPIClient(baseUrl, timeout)
	//
	//	token := apiClient.Login()
	//	dataFromApi := apiClient.GetProducts(token, 0)
	//
	//	data = dataFromApi
	//}

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
