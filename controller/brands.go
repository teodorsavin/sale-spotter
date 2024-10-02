package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"teodorsavin/ah-bonus/database"
)

// GetBrands responds with the list of all brands as JSON.
func GetBrands(c *gin.Context) {
	data := database.AllBrands()

	c.IndentedJSON(http.StatusOK, data)
}
