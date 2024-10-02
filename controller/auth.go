package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
}

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

const secretKey = "randomString"

func getToken(username string) (string, error) {
	hmacSampleSecret := []byte(secretKey)
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
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
		hmacSampleSecret := []byte(secretKey)

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

// Authenticate is the middleware function for JWT authentication.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
				Message: "Unauthorized",
			})
			return
		}

		c.Next()
	}
}
