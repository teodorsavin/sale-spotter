package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

var secretKey = []byte("your-secret-key")

// Claims represents the JWT token claims.
type Claims struct {
	Username string `json:"username"`
	jwt.Token
}

func (c Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c Claims) GetNotBefore() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c Claims) GetIssuer() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c Claims) GetSubject() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c Claims) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
}

// generateToken generates a new JWT token.
func generateToken(username string) (string, error) {
	var t *jwt.Token
	t = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"Username": username,
		})
	return t.SignedString(secretKey)
}

func extractClaims(_ http.ResponseWriter, request *http.Request) (string, error) {
	if request.Header["Token"] != nil {
		tokenString := request.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			return "Error Parsing Token: ", err
		}

		return token.Raw, nil
	}

	return "", nil
}

// getProducts responds with the list of all products as JSON.
func getProducts(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	data := controller.GetAllProducts()

	// If we don't find products in our database, we call the API and save products in our database
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
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	data := controller.AllBrands()

	c.IndentedJSON(http.StatusOK, data)
}

func saveProductsDebug(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

// authenticate is the middleware function for JWT authentication.
func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		log.Println("-----------------------")
		log.Println(tokenString)
		log.Println("-----------------------")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		log.Println("-----------token------------")
		log.Println(token)
		log.Println("-----------token------------")

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Println("-----------claims------------")
		log.Println(claims)
		log.Println("-----------claims------------")

		// Pass the authenticated user to the next handler
		c.Set("username", claims.Username)
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// Routes requiring authentication
	auth := router.Group("/")
	auth.Use(authenticate())
	{
		auth.GET("/brands", getBrands)
		auth.GET("/products", getProducts)
		auth.GET("/test/saveproducts", saveProductsDebug)
	}

	// Token generation endpoint
	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")

		token, err := generateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	router.Run("0.0.0.0:8080")
}
