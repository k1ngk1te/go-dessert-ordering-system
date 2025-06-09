package services

import (
	"dessert-ordering-go-system/models"
)

type CartItemService struct {
	CartItemModel *models.CartItemModel
}

func NewCartItemService(cartItemModel *models.CartItemModel) *CartItemService {
	return &CartItemService{
		CartItemModel: cartItemModel,
	}
}

func (ci *CartItemService) GetCart() (models.Cart, error) {
	return ci.CartItemModel.GetCartItems(1)
}

func (ci *CartItemService) AddCartItem(productID int) error {
	return ci.CartItemModel.AddCartItem(1, productID, 1)
}

func (ci *CartItemService) RemoveSingleCartItem(productID int) error {
	return ci.CartItemModel.RemoveSingleCartItem(1, productID)
}

func (ci *CartItemService) RemoveCartItem(cartItemID int) error {
	return ci.CartItemModel.RemoveCartItem(1, cartItemID)
}

func (ci *CartItemService) Checkout() error {
	return ci.CartItemModel.ClearCart(1)
}
