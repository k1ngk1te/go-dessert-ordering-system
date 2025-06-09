package app

import (
	"database/sql"
	models "dessert-ordering-go-system/models"
)

type ApplicationModels struct {
	// Reference the types from the 'models' package
	CartItem     *models.CartItemModel
	Product      *models.ProductModel
	ProductImage *models.ProductImageModel
	User         *models.UserModel
}

func NewApplicationModels(db *sql.DB) *ApplicationModels {
	return &ApplicationModels{
		CartItem:     &models.CartItemModel{DB: db},
		Product:      &models.ProductModel{DB: db},
		ProductImage: &models.ProductImageModel{DB: db},
		User:         &models.UserModel{DB: db},
	}
}
