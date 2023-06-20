package controller

import (
	"log"

	"teodorsavin/ah-bonus/config"
	"teodorsavin/ah-bonus/model"
)

const (
	selectProductsQuery = `SELECT 
			webshop_id, hq_id, title, sales_unit_size, current_price, price_before_bonus, order_availability_status, 
			main_category, sub_category, brand, available_online, description_highlights, description_full, is_bonus 
		FROM products
		WHERE inserted_at >= DATE_SUB(CURRENT_DATE, INTERVAL WEEKDAY(CURRENT_DATE) DAY) + INTERVAL 1 DAY + INTERVAL '00:00:00' HOUR_SECOND`

	selectBrandsQuery = `SELECT DISTINCT brand
		FROM products
		WHERE inserted_at >= DATE_SUB(CURRENT_DATE, INTERVAL WEEKDAY(CURRENT_DATE) DAY) + INTERVAL 1 DAY + INTERVAL '00:00:00' HOUR_SECOND
		ORDER BY brand ASC`

	insertProductQuery = `INSERT INTO products 
		(webshop_id, hq_id, title, sales_unit_size, current_price, price_before_bonus, order_availability_status,
		main_category, sub_category, brand, available_online, description_highlights, description_full, is_bonus) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

func GetAllProducts() model.BonusProducts {
	var product model.Product
	bonusProducts := model.BonusProducts{}

	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query(selectProductsQuery)
	if err != nil {
		log.Print(err)
		return bonusProducts
	}

	for rows.Next() {
		err = rows.Scan(
			&product.WebshopId,
			&product.HqId,
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
			log.Print(err.Error())
		} else {
			bonusProducts.Products = append(bonusProducts.Products, product)
		}
	}

	return bonusProducts
}

func AllBrands() model.Brands {
	var brand string
	brands := model.Brands{}

	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query(selectBrandsQuery)
	if err != nil {
		log.Print(err)
		return brands
	}

	for rows.Next() {
		err = rows.Scan(&brand)
		if err != nil {
			log.Print(err.Error())
		} else {
			brands = append(brands, brand)
		}
	}

	return brands
}

func InsertProduct(product model.Product) error {
	db := config.ConnectDB()
	defer db.Close()

	_, err := db.Exec(insertProductQuery,
		product.WebshopId, product.HqId, product.Title, product.SalesUnitSize, product.CurrentPrice,
		product.PriceBeforeBonus, product.OrderAvailabilityStatus, product.MainCategory, product.SubCategory,
		product.Brand, product.AvailableOnline, product.DescriptionHighlights, product.DescriptionFull, product.IsBonus)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func InsertProductsBulk(products []model.Product) error {
	db := config.ConnectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	stmt, err := tx.Prepare(insertProductQuery)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, product := range products {
		_, err := stmt.Exec(
			product.WebshopId, product.HqId, product.Title, product.SalesUnitSize, product.CurrentPrice,
			product.PriceBeforeBonus, product.OrderAvailabilityStatus, product.MainCategory, product.SubCategory,
			product.Brand, product.AvailableOnline, product.DescriptionHighlights, product.DescriptionFull, product.IsBonus,
		)
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
