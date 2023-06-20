package model

type Product struct {
	WebshopId               int32           `form:"webshop_id" json:"webshopId"`
	HqId                    int32           `form:"hq_id" json:"hqId"`
	Title                   string          `form:"title" json:"title"`
	SalesUnitSize           string          `form:"sales_unit_size" json:"salesUnitSize"`
	Images                  []Image         `form:"images" json:"images"`
	CurrentPrice            float32         `form:"current_price" json:"currentPrice"`
	PriceBeforeBonus        float32         `form:"price_before_bonus" json:"priceBeforeBonus"`
	OrderAvailabilityStatus string          `form:"order_availability_status" json:"orderAvailabilityStatus"`
	MainCategory            string          `form:"main_category" json:"mainCategory"`
	SubCategory             string          `form:"sub_category" json:"subCategory"`
	Brand                   string          `form:"brand" json:"brand"`
	AvailableOnline         bool            `form:"available_online" json:"availableOnline"`
	DescriptionHighlights   string          `form:"description_highlights" json:"descriptionHighlights"`
	DescriptionFull         string          `form:"description_full" json:"descriptionFull"`
	IsBonus                 bool            `form:"is_bonus" json:"isBonus"`
	DiscountLabels          []DiscountLabel `form:"discount_labels" json:"discountLabels"`
}

type Image struct {
	Width  int32  `form:"width" json:"width"`
	Height int32  `form:"height" json:"height"`
	Url    string `form:"url" json:"url"`
}

type DiscountLabel struct {
	Code               string  `form:"code" json:"code"`
	DefaultDescription string  `form:"default_description" json:"defaultDescription"`
	Price              float32 `form:"price" json:"price"`
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
