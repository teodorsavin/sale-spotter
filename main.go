package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"teodorsavin/ah-bonus/controller"
	"teodorsavin/ah-bonus/model"
	"teodorsavin/ah-bonus/service"
)

const (
	baseUrl = "https://api.ah.nl"
	timeout = time.Second * 30
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
}

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

func getToken(username string) (string, error) {
	hmacSampleSecret := []byte("randomString")
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(hmacSampleSecret)
}

func jwtTokenCheck(tokenString string) (string, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		hmacSampleSecret := []byte("randomString")

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		str := fmt.Sprintf("%v", claims["username"])
		return str, nil
	}

	return "", err
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

// login responds with the token used for other requests
func login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := getToken(form.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, token)
}

func middlewareTokenCheck(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, UnsignedResponse{
			Message: "No Authorization header provided",
		})
		return
	}

	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}
	_, err = jwtTokenCheck(jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	c.Next()
}

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

	router.POST("/login", login)

	privateRouter := router.Group("/api")
	privateRouter.Use(middlewareTokenCheck)
	privateRouter.GET("/products", getProducts)
	privateRouter.GET("/brands", getBrands)

	// to save all products from test.json into the database
	router.GET("/test/saveproducts", saveProductsDebug)

	router.Run("0.0.0.0:8080")
}
