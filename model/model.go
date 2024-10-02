package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	WebshopId               int32           `json:"webshopId"`
	HqId                    int32           `json:"hqId"`
	Title                   string          `json:"title"`
	SalesUnitSize           string          `json:"salesUnitSize"`
	Images                  []Image         `gorm:"foreignKey:ProductID"`
	CurrentPrice            float32         `json:"currentPrice"`
	PriceBeforeBonus        float32         `json:"priceBeforeBonus"`
	OrderAvailabilityStatus string          `json:"orderAvailabilityStatus"`
	MainCategory            string          `json:"mainCategory"`
	SubCategory             string          `json:"subCategory"`
	Brand                   string          `json:"brand"`
	AvailableOnline         bool            `json:"availableOnline"`
	DescriptionHighlights   string          `json:"descriptionHighlights"`
	DescriptionFull         string          `json:"descriptionFull"`
	IsBonus                 bool            `json:"isBonus"`
	DiscountLabels          []DiscountLabel `gorm:"foreignKey:ProductID"`
}

type Image struct {
	gorm.Model
	ProductID uint   `json:"-"` // Foreign key to Product
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
	Url       string `json:"url"`
}

type DiscountLabel struct {
	gorm.Model
	ProductID          uint    `json:"-"` // Foreign key to Product
	Code               string  `json:"code"`
	DefaultDescription string  `json:"defaultDescription"`
	Price              float32 `json:"price"`
}

type BonusProducts struct {
	Products   []Product
	Categories Categories
	Brands     Brands
}

type Brands []string
type Categories []string

type Response struct {
	Data []BonusProducts
}

func (b Brands) ContainsBrand(e string) bool {
	for _, a := range b {
		if a == e {
			return true
		}
	}
	return false
}

func (c Categories) ContainsCategory(e string) bool {
	for _, a := range c {
		if a == e {
			return true
		}
	}
	return false
}
