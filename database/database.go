package database

import (
	"log"

	"github.com/microcosm-cc/bluemonday"

	"teodorsavin/ah-bonus/config"
	"teodorsavin/ah-bonus/model"
)

var db = config.ConnectDB()

const (
	selectProductsQuery = `
		SELECT p.webshop_id, p.hq_id, p.title, p.sales_unit_size, p.current_price, p.price_before_bonus, p.order_availability_status, 
			p.main_category, p.sub_category, p.brand, p.available_online, p.description_highlights, p.description_full, p.is_bonus, i.width, i.height, i.url
		FROM products p
		INNER JOIN images i ON p.id = i.product_id
		WHERE p.inserted_at >= DATE_SUB(CURRENT_DATE, INTERVAL WEEKDAY(CURRENT_DATE) DAY) + INTERVAL 1 DAY + INTERVAL '00:00:00' HOUR_SECOND`

	selectBrandsQuery = `
		SELECT DISTINCT brand
		FROM products
		WHERE inserted_at >= DATE_SUB(CURRENT_DATE, INTERVAL WEEKDAY(CURRENT_DATE) DAY) + INTERVAL 1 DAY + INTERVAL '00:00:00' HOUR_SECOND
		ORDER BY brand ASC`

	insertProductQuery = `
		INSERT INTO products 
		(webshop_id,
		 hq_id,
		 title,
		 sales_unit_size,
		 current_price,
		 price_before_bonus,
		 order_availability_status,
		 main_category,
		 sub_category,
		 brand,
		 available_online,
		 description_highlights,
		 description_full,
		 is_bonus) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	insertImageQuery = `
		INSERT INTO images 
		(product_id, width, height, url) 
		VALUES (?, ?, ?, ?)`
)

func DBQuery(query string, args ...interface{}) (rows model.DBRows, err error) {
	rows, err = db.Query(query, args...)
	if err != nil {
		log.Print(err)
	}
	return rows, err
}

func GetAllProducts() model.BonusProducts {
	bonusProducts := model.BonusProducts{
		Products: make([]model.Product, 0),
	}

	rows, err := DBQuery(selectProductsQuery)
	if err != nil {
		return bonusProducts
	}

	productMap := make(map[int32]*model.Product)

	for rows.Next() {
		var (
			webshopID int32
			image     model.Image
			product   model.Product
		)

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
			&image.Width,
			&image.Height,
			&image.Url,
		)
		if err != nil {
			log.Print(err.Error())
			continue
		}

		webshopID = product.WebshopId
		if existingProduct, ok := productMap[webshopID]; ok {
			existingProduct.Images = append(existingProduct.Images, image)
			productMap[webshopID] = existingProduct
		} else {
			product.Images = []model.Image{image}
			productMap[webshopID] = &product
			bonusProducts.Products = append(bonusProducts.Products, product)
		}
	}

	return bonusProducts
}

func AllBrands() model.Brands {
	var brand string
	brands := model.Brands{}

	rows, err := DBQuery(selectBrandsQuery)
	if err != nil {
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

// StripHTMLTags strips all HTML tags from the input string
func StripHTMLTags(input string) string {
	policy := bluemonday.StrictPolicy()
	// This will strip all HTML tags and return plain text
	return policy.Sanitize(input)
}

func InsertProduct(product model.Product) error {
	descriptionHighlights := StripHTMLTags(product.DescriptionHighlights)
	descriptionFull := StripHTMLTags(product.DescriptionFull)

	// Insert the product
	res, err := db.Exec(insertProductQuery,
		product.WebshopId,
		product.HqId,
		product.Title,
		product.SalesUnitSize,
		product.CurrentPrice,
		product.PriceBeforeBonus,
		product.OrderAvailabilityStatus,
		product.MainCategory,
		product.SubCategory,
		product.Brand,
		product.AvailableOnline,
		descriptionHighlights,
		descriptionFull,
		product.IsBonus)
	if err != nil {
		log.Print(err)
		return err
	}

	// Get the auto-incremented product_id
	productID, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return err
	}

	// Insert the images associated with the product
	for _, image := range product.Images {
		_, err := db.Exec(insertImageQuery,
			productID, image.Width, image.Height, image.Url)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

func InsertProductsBulk(products []model.Product) error {
	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	productStmt, err := tx.Prepare(insertProductQuery)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return err
	}
	defer productStmt.Close()

	imageStmt, err := tx.Prepare(insertImageQuery)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return err
	}
	defer imageStmt.Close()

	for _, product := range products {
		descriptionHighlights := StripHTMLTags(product.DescriptionHighlights)
		descriptionFull := StripHTMLTags(product.DescriptionFull)

		// Insert the product
		res, err := productStmt.Exec(
			product.WebshopId,
			product.HqId,
			product.Title,
			product.SalesUnitSize,
			product.CurrentPrice,
			product.PriceBeforeBonus,
			product.OrderAvailabilityStatus,
			product.MainCategory,
			product.SubCategory,
			product.Brand,
			product.AvailableOnline,
			descriptionHighlights,
			descriptionFull,
			product.IsBonus,
		)
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}

		// Get the auto-incremented product_id
		productID, err := res.LastInsertId()
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}

		// Insert the images associated with the product
		for _, image := range product.Images {
			_, err := imageStmt.Exec(
				productID, image.Width, image.Height, image.Url,
			)
			if err != nil {
				log.Print(err)
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
