package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"teodorsavin/ah-bonus/model"
)

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

func (c *APIClient) GetProducts(accessToken string, page int32) model.BonusProducts {
	pageNumber := int32(0)
	totalPages := int32(0)

	req := c.BuildGetProductsRequest(accessToken, page)
	res := c.DoRequest(req)
	searchData := c.ReadResponseBodyGetProducts(res)

	pageNumber = searchData.Page.Number
	totalPages = searchData.Page.TotalPages
	bonusProducts := c.GetBonusProducts(searchData)

	if pageNumber < totalPages {
		pageNumber++
		sleepDuration := 2 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), sleepDuration)
		defer cancel()

		select {
		case <-time.After(sleepDuration):
			nextProducts := c.GetProducts(accessToken, pageNumber)
			bonusProducts.Products = append(bonusProducts.Products, nextProducts.Products...)
		case <-ctx.Done():
			log.Println("Sleep duration exceeded, continuing without the next products.")
		}
	}

	return bonusProducts
}

func (c *APIClient) BuildGetProductsRequest(accessToken string, page int32) *http.Request {
	baseURL := c.BaseURL + "/mobile-services/product/search/v2?query=Drogisterij"
	if page > 0 {
		baseURL = baseURL + "&page=" + strconv.FormatInt(int64(page), 10)
	}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatalf("impossible to build GetProducts request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Application", "AHWEBSHOP")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return req
}

func (c *APIClient) ReadResponseBodyGetProducts(response *http.Response) *ResponseSearch {
	defer response.Body.Close()

	searchData := &ResponseSearch{}
	err := json.NewDecoder(response.Body).Decode(searchData)
	if err != nil {
		log.Panic(err)
	}
	return searchData
}

func (c *APIClient) GetBonusProducts(data *ResponseSearch) model.BonusProducts {
	bonusProducts := model.BonusProducts{}
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
