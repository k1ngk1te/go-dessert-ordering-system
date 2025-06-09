package services

import (
	appErrors "dessert-ordering-go-system/internal/app_errors"
	models "dessert-ordering-go-system/models"
)

type ProductService struct {
	ProductModel *models.ProductModel
}

func NewProductService(productModel *models.ProductModel) *ProductService {
	return &ProductService{
		ProductModel: productModel,
	}
}

func (ps *ProductService) GetAllProducts() ([]*models.Product, error) {
	return ps.ProductModel.GetAllProducts()
}

func (ps *ProductService) GetProductDetail(productID int) (*models.Product, error) {
	product, err := ps.ProductModel.GetProductByID(productID)

	if err == models.ErrProductNotFound || product == nil {
		return nil, &appErrors.NotFoundError{Message: err.Error(), Code: 404}
	} else if err != nil {
		return nil, err
	}

	return product, nil
}
