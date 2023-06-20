package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	controller "teodorsavin/ah-bonus/controller"
	model "teodorsavin/ah-bonus/model"
)

type LoginData struct {
	ClientId string `json:"clientId"`
}

type AuthData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

type ResponseSearch struct {
	Page     Page            `json:"page"`
	Products []model.Product `json:"products"`
}

type Page struct {
	Size          int32 `json:"size"`
	TotalElements int32 `json:"totalElements"`
	TotalPages    int32 `json:"totalPages"`
	Number        int32 `json:"number"`
}

func Login() string {
	loginData := LoginData{
		ClientId: "appie",
	}
	// marshall data to json (like json_encode)
	marshalled, err := json.Marshal(loginData)
	if err != nil {
		log.Fatalf("impossible to marshall loginData: %s", err)
	}

	req, err := http.NewRequest("POST", "https://api.ah.nl/mobile-auth/v1/auth/token/anonymous", bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")

	res := DoRequest(req)
	accessToken := ReadResponseBodyLogin(res)
	return accessToken
}

func DoRequest(request *http.Request) *http.Response {
	// create http client
	// do not forget to set timeout; otherwise, no timeout!
	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	res, err := client.Do(request)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)

	return res
}

func ReadResponseBodyLogin(response *http.Response) (accessToken string) {
	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer response.Body.Close()

	authData := &AuthData{}
	err := json.NewDecoder(response.Body).Decode(authData)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != http.StatusOK {
		panic(response.Status)
	}

	accessToken = authData.AccessToken
	return accessToken
}

func GetProducts(accessToken string, page int32) (bonusProducts model.BonusProducts) {
	pageNumber := int32(0)
	totalPages := int32(0)

	req := BuildGetProductsRequest(accessToken, page)
	res := DoRequest(req)
	searchData := ReadResponseBodyGetProducts(res)

	pageNumber = searchData.Page.Number
	totalPages = searchData.Page.TotalPages
	bonusProducts = GetBonusProducts(searchData)

	if pageNumber < totalPages {
		pageNumber++
		// Wait 2 seconds before doing a new call
		time.Sleep(2 * time.Second)
		nextProducts := GetProducts(accessToken, pageNumber)
		bonusProducts.Products = append(bonusProducts.Products, nextProducts.Products...)
	}

	return bonusProducts
}

func BuildGetProductsRequest(accessToken string, page int32) (request *http.Request) {
	baseUrl := "https://api.ah.nl/mobile-services/product/search/v2?query=Drogisterij"
	if page > 0 {
		baseUrl = baseUrl + "&page=" + strconv.Itoa(int(page))
	}
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		log.Fatalf("impossible to build GetProducts request: %s", err)
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Application", "AHWEBSHOP")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return req
}

func ReadResponseBodyGetProducts(response *http.Response) (searchData *ResponseSearch) {
	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer response.Body.Close()

	searchData = &ResponseSearch{}
	err := json.NewDecoder(response.Body).Decode(searchData)
	if err != nil {
		panic(err)
	}
	return searchData
}

func GetBonusProducts(data *ResponseSearch) (bonusProducts model.BonusProducts) {
	for _, item := range data.Products {
		if item.IsBonus {
			bonusProducts.Products = append(bonusProducts.Products, item)
			if !bonusProducts.Brands.ContainsBrand(item.Brand) {
				bonusProducts.Brands = append(bonusProducts.Brands, item.Brand)
			}
			if !bonusProducts.Categories.ContainsCategory(item.SubCategory) {
				bonusProducts.Categories = append(bonusProducts.Categories, item.SubCategory)
			}
		}
	}

	return bonusProducts
}

// getProducts responds with the list of all products as JSON.
func getProducts(c *gin.Context) {
	data := controller.AllProducts()

	//token := Login()
	//data := GetProducts(token, 0)

	c.IndentedJSON(http.StatusOK, data)
}

// getProducts responds with the list of all brands as JSON.
func getBrands(c *gin.Context) {
	data := controller.AllBrands()

	c.IndentedJSON(http.StatusOK, data)
}

func saveProductsDebug(c *gin.Context) {
	filename := "../usr/src/app/database-sample/test.json"

	// Get the absolute path of the file based on the current working directory
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return
	}

	if _, err = os.Stat(absPath); err == nil {
		fmt.Printf("File exists\n")
	} else {
		fmt.Printf("File does not exist\n")
		fmt.Printf(absPath + "\n")
	}

	file, _ := os.ReadFile(filename)
	var data model.BonusProducts
	err = json.Unmarshal(file, &data)
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
