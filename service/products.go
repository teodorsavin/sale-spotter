package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	totalPages := int32(0)
	productsCh := make(chan []model.Product)
	bonusProducts := model.BonusProducts{}

	req := c.BuildGetProductsRequest(accessToken, page)
	res := c.DoRequest(req)
	searchData := c.ReadResponseBodyGetProducts(res)

	totalPages = searchData.Page.TotalPages
	bonusProducts = c.GetBonusProducts(searchData)

	// Fetch pages concurrently
	for i := int32(1); i <= totalPages; i++ {
		go func(page int32) {
			nextProducts := c.GetProductsForPage(accessToken, page)
			productsCh <- nextProducts.Products
		}(i)
	}

	// Collect all products from the channel
	for i := int32(1); i <= totalPages; i++ {
		products := <-productsCh
		bonusProducts.Products = append(bonusProducts.Products, products...)
	}

	return bonusProducts
}

func (c *APIClient) GetProductsForPage(accessToken string, page int32) model.BonusProducts {
	req := c.BuildGetProductsRequest(accessToken, page)
	res := c.DoRequest(req)
	searchData := c.ReadResponseBodyGetProducts(res)
	return c.GetBonusProducts(searchData)
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
