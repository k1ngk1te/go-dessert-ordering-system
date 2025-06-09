package models

import (
	"database/sql"
	"log"
	"sort"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Thumbnail   string    `json:"thumbnail"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProductImage struct {
	ID        int       `json:"id"`
	ProductID int       `json:"productId"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProductModel struct {
	DB *sql.DB
}

type ProductImageModel struct {
	DB *sql.DB
}

func (m *ProductModel) GetAllProducts() ([]*Product, error) {
	// 1. Construct the SQL query.
	query := `
					SELECT
						p.id,
						p.title,
						p.category,
						p.description,
						p.price,
						p.thumbnail,
						p.created_at,
						p.updated_at,
						pi.image
					FROM
						products AS p
					LEFT JOIN
						product_images AS pi 
					ON 
						p.id = pi.product_id
					ORDER BY
						p.id, pi.id
				
	`

	// 2. Execute the query. db.Query() returns a *sql.Rows (multiple rows)
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Printf("ERROR: ProductModel.GetAllProducts - m.DB.Query: %v", err)
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after we're done

	// 3. Iterate over the results and scan them into Product structs.
	productsMap := make(map[int]*Product)

	for rows.Next() {
		var (
			productID   int
			title       string
			category    string
			description string
			price       float64
			thumbnail   string
			createdAt   time.Time
			updatedAt   time.Time
			imageURL    sql.NullString // Use sql.NullString to handle NULL image_url values from LEFT JOIN
		)

		err := rows.Scan(
			&productID,
			&title,
			&category,
			&description,
			&price,
			&thumbnail,
			&createdAt,
			&updatedAt,
			&imageURL, // Scan into sql.NullString
		)
		if err != nil {
			log.Printf("ERROR: ProductModel.GetAllProducts - rows.Scan: %v", err)
			return nil, err
		}

		// Check if the product is already in our map (i.e., we've seen its details before)
		p, ok := productsMap[productID]
		if !ok {
			// If not, this is a new product, so create it and add it to the map
			p = &Product{
				ID:          productID,
				Title:       title,
				Category:    category,
				Description: description,
				Price:       price,
				Thumbnail:   thumbnail,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Images:      []string{},
			}
			productsMap[productID] = p
		}

		// IF the imageURL is not NULL (meaning there's an image for this product)
		if imageURL.Valid {
			p.Images = append(p.Images, imageURL.String)
		}
	}

	// 4. Check for errors from iterating over the rows.
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: ProductModel.GetAllProducts - rows.Err: %v", err)
		return nil, err
	}

	products := []*Product{}
	for _, p := range productsMap {
		products = append(products, p)
	}

	// Optional: Sort products by ID or another field if you want a consistent order
	sort.Slice(products, func(i, j int) bool {
		return products[i].ID < products[j].ID
	})

	return products, nil
}

func (m *ProductModel) GetProductByID(ProductID int) (*Product, error) {
	// 1. Query
	query := `
		SELECT
			p.id,
			p.title,
			p.category,
			p.description,
			p.price,
			p.thumbnail,
			p.created_at,
			p.updated_at,
			pi.image
		FROM
			products AS p	
		LEFT JOIN
			product_images AS pi
		ON 
			pi.product_id = p.id
		WHERE
			p.id = ?
		ORDER BY
			pi.id
		
	`

	// 2. Run the query
	rows, err := m.DB.Query(query, ProductID)
	if err != nil {
		log.Printf("ERROR: ProductModel.GetProductByID - m.DB.Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 3. create a single product struct to take in the data
	var product *Product

	for rows.Next() {
		var (
			productID   int
			title       string
			category    string
			description string
			price       float64
			thumbnail   string
			createdAt   time.Time
			updatedAt   time.Time
			imageURL    sql.NullString
		)

		err := rows.Scan(
			&productID,
			&title,
			&category,
			&description,
			&price,
			&thumbnail,
			&createdAt,
			&updatedAt,
			&imageURL,
		)
		if err != nil {
			log.Printf("ERROR: ProductModel.GetProductByID - rows.Scan: %v", err)
			return nil, err
		}

		// Check if the product has been set or not
		if product == nil {
			product = &Product{
				ID:          productID,
				Title:       title,
				Category:    category,
				Description: description,
				Price:       price,
				Thumbnail:   thumbnail,
				Images:      []string{},
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}
		}

		// Check if the image is valid
		if imageURL.Valid {
			product.Images = append(product.Images, imageURL.String)
		}
	}

	// Check for rows error
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: ProductModel.GetProductByID - rows.Err: %v", err)
		return nil, err
	}

	// Check no rows found
	if product == nil {
		return nil, ErrProductNotFound
	}

	return product, nil
}
