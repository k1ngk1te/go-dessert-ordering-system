package models

import (
	"database/sql"
	"log"
	"time"
)

type Cart = []*CartItem

// CartItem represents a single item in a user's cart in the database.
type CartItem struct {
	ID        int       `json:"id"`
	ProductID int       `json:"productId"`
	Quantity  int       `json:"quantity"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Product *CartItemProduct
}

// CartItemSimplified is a stripped-down version for display purposes.
type CartItemSimplified struct {
	ID        int
	UserID    int
	ProductID int
	Quantity  int
}

// CartDisplayItem combines product details with cart item details for rendering.
type CartDisplayItem struct {
	Product    *CartItemProduct
	CartItem   *CartItemSimplified
	TotalPrice float64
}

type CartItemProduct struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Thumbnail   string  `json:"thumbnail"`
}

type CartItemModel struct {
	DB *sql.DB
}

func (m *CartItemModel) GetCartItems(userID int) ([]*CartItem, error) {
	// 1. Construct the SQL query.
	query := `
					SELECT
						ci.id,
						ci.product_id,
						ci.quantity,
						ci.user_id,
						ci.created_at,
						ci.updated_at,

						p.id,
						p.title,
						p.category,
						p.description,
						p.price,
						p.thumbnail
					FROM
						cart_items as ci
					LEFT JOIN
						products AS p
					ON
						ci.product_id = p.id
					WHERE
						ci.user_id = ?
					ORDER BY ci.id DESC
	`

	// 2. Execute the query
	rows, err := m.DB.Query(query, userID)
	if err != nil {
		log.Printf("ERROR: CartItemModel.GetCartItems - m.DB.Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 3. Iterate over the results and scan them into cart items structs
	cartItems := make([]*CartItem, 0)

	for rows.Next() {
		cartItem := CartItem{}

		// Initialize the Product sub-struct BEFORE scanning into it!
		cartItem.Product = &CartItemProduct{} // <--- CRITICAL FIX

		err := rows.Scan(
			&cartItem.ID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&cartItem.UserID,
			&cartItem.CreatedAt,
			&cartItem.UpdatedAt,

			&cartItem.Product.ID,
			&cartItem.Product.Title,
			&cartItem.Product.Category,
			&cartItem.Product.Description,
			&cartItem.Product.Price,
			&cartItem.Product.Thumbnail,
		)
		if err != nil {
			log.Printf("ERROR: CartItemModel.GetCartItems - rows.Scan: %v", err)
			return nil, err
		}
		cartItems = append(cartItems, &cartItem)
	}

	// 4. Check for errors from iterating over the rows.
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: CartItem.GetCartItems - rows.Err: %v", err)
		return nil, err
	}

	return cartItems, nil
}

func (m *CartItemModel) AddCartItem(userID, productID, quantity int) error {
	// -1. Start a database transaction
	// 0. Check if the product is on the database
	// 1. Check if the product is already in the user cart
	// 2. If it's not perform an insert
	// 3. If it's there perform an update (and skip No. 3)

	// -1
	tx, err := m.DB.Begin()
	if err != nil {
		log.Printf("ERROR: CartItemModel.AddCartItem - m.DB.Begin: %v", err)
		return err
	}

	defer tx.Rollback()

	// 0
	productQuery := `
				SELECT
					p.id
				FROM 
					products AS p
				WHERE
					p.id = ?
	`

	var pID int

	productRow := tx.QueryRow(productQuery, productID)
	err = productRow.Scan(&pID)

	if err == sql.ErrNoRows {
		log.Printf("ERROR: CartItemModel.AddCartItem - Product with ID %d not found.", productID)
		return ErrProductNotFound
	} else if err != nil {
		log.Printf("ERROR: CartItemModel.AddCartItem - productRow.Scan: %v", err)
		return err
	}

	// 1.
	cartItem := &CartItemSimplified{}

	cartItemQuery := `
				SELECT
					ci.id,
					ci.user_id,
					ci.product_id,
					ci.quantity
				FROM
					cart_items as ci
				WHERE
					ci.product_id = ? AND ci.user_id = ?
				LIMIT 1
	`

	cartRow := tx.QueryRow(cartItemQuery, productID, userID)
	err = cartRow.Scan(
		&cartItem.ID,
		&cartItem.UserID,
		&cartItem.ProductID,
		&cartItem.Quantity,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ERROR: CartItemModel.AddCartItem - cartRow.Scan: %v", err)
		return err
	}

	// 2.
	if err == sql.ErrNoRows {
		insertQuery := `
		INSERT INTO cart_items(
			product_id, user_id, quantity, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?)
	`
		_, err := tx.Exec(insertQuery, productID, userID, quantity, time.Now(), time.Now())
		if err != nil {
			log.Printf("ERROR: CartItemModel.AddCartItem - Create CartItem - m.DB.Exec: %v", err)
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.Printf("ERROR: CartItemModel.AddCartItem - tx.Commit failed: %v", err)
			return err
		}

		return nil
	} else if err != nil {
		log.Printf("ERROR: CartItem.AddCartItem - cartRow.Scan: %v", err)
		return err
	}

	// 3
	updateQuery := `
			UPDATE cart_items
			SET
				cart_items.quantity = ?,
				cart_items.updated_at = ?
			WHERE 
				cart_items.id = ?
		`
	_, err = tx.Exec(updateQuery, cartItem.Quantity+quantity, time.Now(), cartItem.ID)
	if err != nil {
		log.Printf("ERROR: CartItemModel.AddCartItem - Update CartItem - m.DB.Exec: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("ERROR: CartItemModel.AddCartItem - tx.Commit failed: %v", err)
		return err
	}

	return nil
}

func (m *CartItemModel) RemoveCartItem(userID int, cartItemID int) error {
	deleteQuery := `
		DELETE FROM cart_items
		WHERE id = ? AND user_id = ?
	`

	result, err := m.DB.Exec(deleteQuery, cartItemID, userID)
	if err != nil {
		log.Printf("ERROR: CartItemModel.RemoveCartItem - m.DB.Exec: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("ERROR: CartItemModel.RemoveCartItem - result.RowsAffected failed for cartItemID %d, userID %d: %v", cartItemID, userID, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("cart item with the specified ID - %d not found.", cartItemID)
		return ErrCartItemNotFound
	}
	return nil
}

func (m *CartItemModel) RemoveSingleCartItem(userID, productID int) error {
	// Step 0: Start a database transaction for atomicity
	tx, err := m.DB.Begin()
	if err != nil {
		log.Printf("ERROR: CartItemModel.RemoveSingleCartItem - m.DB.Begin failed: %v", err)
		return err
	}
	defer tx.Rollback()

	// Step 1: Check if the cart item exists
	// Step 2: If new quantity is less than 1, then delete the item else update

	// 1.
	selectQuery := `
		SELECT 
			ci.id,
			ci.user_id,
			ci.product_id,
			ci.quantity
		FROM
			cart_items AS ci
		WHERE
			user_id = ? AND product_id = ?
		LIMIT 1
	`

	cartItem := CartItemSimplified{}

	row := tx.QueryRow(selectQuery, userID, productID)
	err = row.Scan(
		&cartItem.ID,
		&cartItem.UserID,
		&cartItem.ProductID,
		&cartItem.Quantity,
	)

	if err == sql.ErrNoRows {
		log.Printf("no cart item with the specified userID %d and productID %d exists.", userID, productID)
		return ErrCartItemNotFound
	} else if err != nil {
		log.Printf("ERROR: CartItemModel.RemoveSingleCartItem - row.Scan: %v", err)
		return err
	}

	cartItem.Quantity--

	if cartItem.Quantity < 1 {
		deleteQuery := `
			DELETE FROM cart_items
			WHERE id = ? AND user_id = ?
		`
		_, err := tx.Exec(deleteQuery, cartItem.ID, cartItem.UserID)
		if err != nil {
			log.Printf("ERROR: CartItemModel.RemoveSingleCartItem - Delete Cart Item - m.DB.Exec: %v", err)
			return err
		}
	} else {
		updateQuery := `
			UPDATE cart_items
			SET 
				cart_items.quantity = ?,
				updated_at = ?
			WHERE
				cart_items.id = ? AND user_id = ?
		`
		_, err := tx.Exec(updateQuery, cartItem.Quantity, time.Now(), cartItem.ID, cartItem.UserID)
		if err != nil {
			log.Printf("ERROR: CartItemModel.RemoveSingleCartItem - Update Cart Item - m.DB.Exec: %v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("ERROR: CartItemModel.RemoveSingleCartItem - tx.Commit failed: %v", err)
		return err
	}

	return nil
}

func (m *CartItemModel) ClearCart(userID int) error {
	deleteQuery := `
		DELETE FROM cart_items
		WHERE user_id = ?
	`

	result, err := m.DB.Exec(deleteQuery, userID)
	if err != nil {
		log.Printf("ERROR: CartItemModel.ClearCart = m.DB.Exec: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("ERROR: CartItemModel.ClearCart - result.RowsAffected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("no cart items were found for user ID: %d.", userID)
		return ErrNoCartItemsFound
	}

	return nil
}
