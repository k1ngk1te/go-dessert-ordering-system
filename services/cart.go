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

func (ci *CartItemService) GetCart(userID int) (models.Cart, error) {
	return ci.CartItemModel.GetCartItems(userID)
}

func (ci *CartItemService) AddCartItem(userID, productID int) error {
	return ci.CartItemModel.AddCartItem(userID, productID, 1)
}

func (ci *CartItemService) RemoveSingleCartItem(userID, productID int) error {
	return ci.CartItemModel.RemoveSingleCartItem(userID, productID)
}

func (ci *CartItemService) RemoveCartItem(userID, cartItemID int) error {
	return ci.CartItemModel.RemoveCartItem(userID, cartItemID)
}

func (ci *CartItemService) Checkout(userID int) error {
	return ci.CartItemModel.ClearCart(userID)
}
