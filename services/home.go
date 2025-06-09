package services

import (
	"fmt"
	"log"

	"dessert-ordering-go-system/models"
)

type ApplicationCartItem struct {
	CartItem   *models.CartItem
	Product    *models.Product
	Quantity   int
	TotalPrice float64
}

type ApplicationProduct struct {
	Product  *models.Product
	Quantity int
}

type HomeTemplateData struct {
	Cart              []ApplicationCartItem
	CsrfToken         string
	Errors            []string
	Messages          []string
	Products          []ApplicationProduct
	IsCartEmpty       bool
	TotalCartPrice    float64
	TotalCartQuantity int
}

func (c HomeTemplateData) String() string {
	return fmt.Sprintf("Cart [%v], CsrfToken: %v, Errors [%v], Messages [%v], Products [%v], IsCartEmpty: %v, TotalCartPrice: %v, TotalCartQuantity: %v",
		len(c.Cart), c.CsrfToken,
		len(c.Errors),
		len(c.Messages),
		len(c.Products),
		c.IsCartEmpty,
		c.TotalCartPrice,
		c.TotalCartQuantity)
}

type HomeTemplateDataService struct {
	CartItemModel *models.CartItemModel
	ProductModel  *models.ProductModel
}

type GetHomeTemplateContentOptionsFunc func(*HomeTemplateData)

func NewHomeTemplateDataService(cm *models.CartItemModel, pm *models.ProductModel) *HomeTemplateDataService {
	return &HomeTemplateDataService{
		CartItemModel: cm,
		ProductModel:  pm,
	}
}

func (s *HomeTemplateDataService) WithCsrfToken(csrfToken string) GetHomeTemplateContentOptionsFunc {
	return func(opts *HomeTemplateData) {
		opts.CsrfToken = csrfToken
	}
}

func (s *HomeTemplateDataService) WithErrors(errs []string) GetHomeTemplateContentOptionsFunc {
	return func(opts *HomeTemplateData) {
		opts.Errors = append(opts.Errors, errs...)
	}
}

func (s *HomeTemplateDataService) GetHomeTemplateContent(opts ...GetHomeTemplateContentOptionsFunc) (*HomeTemplateData, error) {

	var templateContent *HomeTemplateData = &HomeTemplateData{}

	for _, fn := range opts {
		fn(templateContent)
	}

	userID := 1
	products, err := s.ProductModel.GetAllProducts()
	if err != nil {
		log.Printf("ERROR: HomeTemplateDataService.GetHomeTemplateContent - Failed to get all products: %v", err)
		return nil, fmt.Errorf("failed to load product catalog: %w", err)
	}

	cart, err := s.CartItemModel.GetCartItems(userID)
	if err != nil {
		log.Printf("ERROR: HomeTemplateDataService.GetHomeTemplateContent - Failed to get user cart items for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to load user cart: %w", err)
	}

	var errors []string

	var applicationProducts []ApplicationProduct = make([]ApplicationProduct, 0, len(products))
	var applicationCartItems []ApplicationCartItem = make([]ApplicationCartItem, 0, len(cart))

	cartQuantities := make(map[int]int)

	totalCartPrice := 0.0 // Initialize total cart price
	totalCartQuantity := 0

	for _, cartItem := range cart {
		cartQuantities[cartItem.ProductID] = cartItem.Quantity

		product, err := s.ProductModel.GetProductByID(cartItem.ProductID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("could not retrieve product with ID: %v, %v", cartItem.ProductID, err))
			continue
		}

		totalPrice := float64(cartItem.Quantity) * product.Price

		totalCartPrice += totalPrice
		totalCartQuantity += cartItem.Quantity

		applicationCartItems = append(applicationCartItems, ApplicationCartItem{
			CartItem:   cartItem,
			Product:    product,
			Quantity:   cartItem.Quantity,
			TotalPrice: totalPrice,
		})
	}

	for index, product := range products {
		quantity := cartQuantities[product.ID]

		applicationProducts = append(applicationProducts, ApplicationProduct{
			Product:  products[index],
			Quantity: quantity,
		})
	}

	templateContent.Cart = applicationCartItems
	templateContent.Errors = errors
	templateContent.Messages = []string{}
	templateContent.IsCartEmpty = totalCartQuantity < 1
	templateContent.Products = applicationProducts
	templateContent.TotalCartPrice = totalCartPrice
	templateContent.TotalCartQuantity = totalCartQuantity

	return templateContent, nil
}
