package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	
	appErrors "dessert-ordering-go-system/internal/app_errors"
	responses "dessert-ordering-go-system/internal/response"
)

func (h *WebHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.Services.Product.GetAllProducts()
	if err != nil {
		response := responses.NewErrorJsonResponse("Failed to fetch products: " + err.Error())
		responses.WriteJsonResponse(w, http.StatusInternalServerError, response)
		return
	}
	response := responses.NewSuccessJsonDataResponse("Fetched Products", products)
	responses.WriteJsonResponse(w, http.StatusOK, response)
}

func (h *WebHandler) GetProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	pathId := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(pathId)

	if err != nil {
		response := responses.NewErrorJsonResponse("Invalid Product ID")
		responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		return
	}

	product, err := h.Services.Product.GetProductDetail(productID)
	if err != nil {
		response := responses.NewErrorJsonResponse(err.Error())
		statusCode := http.StatusInternalServerError
		// Check if it's a not found error
		if _, ok := err.(*appErrors.NotFoundError); ok {
			statusCode = http.StatusNotFound
		}
		responses.WriteJsonResponse(w, statusCode, response)
		return
	}
	response := responses.NewSuccessJsonDataResponse("Fetched Product", product)
	responses.WriteJsonResponse(w, http.StatusOK, response)
}