package app

import (
	services "dessert-ordering-go-system/services"
)

type ApplicationServices struct {
	// Reference the types from the 'services' package
	Auth                 *services.AuthService
	CartItem             *services.CartItemService
	Product              *services.ProductService
	HomeTemplateData     *services.HomeTemplateDataService
	LoginTemplateData    *services.LoginTemplateDataService
	RegisterTemplateData *services.RegisterTemplateDataService
}

func NewApplicationServices(models *ApplicationModels) *ApplicationServices {
	return &ApplicationServices{
		Auth:                 services.NewAuthService(models.User),
		CartItem:             services.NewCartItemService(models.CartItem),
		Product:              services.NewProductService(models.Product),
		HomeTemplateData:     services.NewHomeTemplateDataService(models.CartItem, models.Product),
		LoginTemplateData:    services.NewLoginTemplateDataService(),
		RegisterTemplateData: services.NewRegisterTemplateDataService(),
	}
}
