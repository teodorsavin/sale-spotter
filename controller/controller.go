package controller

import (
	"log"
	config "teodorsavin/ah-bonus/Config"
	model "teodorsavin/ah-bonus/Model"
)

func AllProducts() model.BonusProducts {
	var product model.Product
	bonusProducts := model.BonusProducts{}

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query(
		// images, discountLabels
		`SELECT 
    			webshop_id, hq_id, title, sales_unit_size, current_price, price_before_bonus, order_availability_status, 
    			main_category, sub_category, brand, available_online, description_highlights, description_full, is_bonus 
			FROM products`,
	)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(
			&product.WebshopId,
			&product.HqId,
			&product.Title,
			&product.Title,
			&product.SalesUnitSize,
			&product.CurrentPrice,
			&product.PriceBeforeBonus,
			&product.OrderAvailabilityStatus,
			&product.MainCategory,
			&product.SubCategory,
			&product.Brand,
			&product.AvailableOnline,
			&product.DescriptionHighlights,
			&product.DescriptionFull,
			&product.IsBonus,
		)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			bonusProducts.Products = append(bonusProducts.Products, product)
		}
	}

	return bonusProducts
}

func AllBrands() model.Brands {
	var brand string
	brands := model.Brands{}

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query(
		// images, discountLabels
		`SELECT DISTINCT brand
				FROM products
				WHERE inserted_at >= DATE_SUB(CURRENT_DATE, INTERVAL WEEKDAY(CURRENT_DATE) DAY) + INTERVAL 1 DAY + INTERVAL '00:00:00' HOUR_SECOND
				ORDER BY brand ASC;`,
	)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(&brand)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			brands = append(brands, brand)
		}
	}

	return brands
}

func InsertProduct(product model.Product) {
	db := config.Connect()
	defer db.Close()

	webshopId := product.WebshopId
	hqId := product.HqId
	title := product.Title
	salesUnitSize := product.SalesUnitSize
	currentPrice := product.CurrentPrice
	priceBeforeBonus := product.PriceBeforeBonus
	orderAvailabilityStatus := product.OrderAvailabilityStatus
	mainCategory := product.MainCategory
	subCategory := product.SubCategory
	brand := product.Brand
	availableOnline := product.AvailableOnline
	descriptionHighlights := product.DescriptionHighlights
	descriptionFull := product.DescriptionFull
	isBonus := product.IsBonus

	_, err := db.Exec(`INSERT INTO products 
    	(webshop_id, hq_id, title, sales_unit_size, current_price, price_before_bonus, order_availability_status,
		main_category, sub_category, brand, available_online, description_highlights, description_full, is_bonus) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		webshopId, hqId, title, salesUnitSize, currentPrice, priceBeforeBonus, orderAvailabilityStatus, mainCategory,
		subCategory, brand, availableOnline, descriptionHighlights, descriptionFull, isBonus)
	if err != nil {
		log.Print(err)
		return
	}
}

func InsertProductsBulk(products []model.Product) {
	for _, product := range products {
		InsertProduct(product)
	}
}
