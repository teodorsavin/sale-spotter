package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"teodorsavin/ah-bonus/database"
	"teodorsavin/ah-bonus/model"
)

func SaveProductsDebug(c *gin.Context) {
	filename := "../usr/src/app/database-sample/test.json"

	file, _ := os.ReadFile(filename)
	var data model.BonusProducts
	err := json.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("impossible to read/unmarshall test.json: %s", err)
	}
	err = database.InsertProductsBulk(data.Products)
	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}
