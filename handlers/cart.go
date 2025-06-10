package handlers

import (
	"dessert-ordering-go-system/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	appConstants "dessert-ordering-go-system/internal/app_constants"
	appErrors "dessert-ordering-go-system/internal/app_errors"
	responses "dessert-ordering-go-system/internal/response"

	"github.com/go-chi/chi/v5"
)

type AddOrder struct {
	ProductID int `json:"productId" form:"productId" validate:"required,min=1"`
}

func (h *WebHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")

	csrfToken := h.Session.GetCsrfToken(r.Context())
	userID := h.Session.GetAuthUserID(r.Context())

	if strings.HasPrefix(acceptType, "application/json") {
		cart, err := h.Services.CartItem.GetCart(userID)
		if err != nil {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusInternalServerError, response)
			return
		}
		response := responses.NewSuccessJsonDataResponse("Fetched Cart", cart)
		responses.WriteJsonHeadersResponse(w, http.StatusOK, response, map[string]string{appConstants.X_CSRF_Token: csrfToken})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *WebHandler) AddCartItemHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")
	contentType := r.Header.Get("Content-Type")

	var addOrder AddOrder

	csrfToken := h.Session.GetCsrfToken(r.Context())
	userID := h.Session.GetAuthUserID(r.Context())

	// Json Content
	if strings.HasPrefix(contentType, "application/json") {

		errStatusCode, err := JsonBodyDecoder(w, r, &addOrder)
		if err != nil {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, errStatusCode, response)
			return
		}

		if addOrder.ProductID < 1 {
			response := responses.NewErrorJsonResponse("product ID is required.")
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
			return
		}
	} else {
		// Html Content
		err := r.ParseForm()
		if err != nil {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: AddCartItemHandler - ParseForm - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, fmt.Sprintf("ERROR: Failed to parse form data: %v", err))
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusInternalServerError)
			return
		}

		// Get the productId from the form and convert to an int
		productIdStr := r.FormValue("productId")
		productId, err := strconv.Atoi(productIdStr)
		if err != nil {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: AddCartItemHandler - Get Product - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, fmt.Sprintf("ERROR: Invalid product ID received: %s - %v", productIdStr, err))
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusBadRequest)
			return
		}

		addOrder.ProductID = productId
	}

	// -- Validation Goes Here --
	// -- Perform Validation --
	validationErrors := h.Validator.ValidateStruct(addOrder)
	if validationErrors != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonDataResponse("Validation failed", validationErrors)
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostRegisterHandler - GetRegisterTemplateContent on validation error: %v", templateDataErr)
				http.Error(w, "Failed to reload index page after validation error", http.StatusInternalServerError)
				return
			}

			for field, msg := range validationErrors {
				data.Errors = append(data.Errors, fmt.Sprintf("%s: %s", field, msg))
			}
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusBadRequest) // 400 Bad Request
		}
		return
	}

	err := h.Services.CartItem.AddCartItem(userID, addOrder.ProductID)
	if err != nil {
		if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			// Product not found when trying to add to cart
			if strings.HasPrefix(acceptType, "application/json") {
				response := responses.NewErrorJsonResponse(notFoundErr.Error())
				responses.WriteJsonResponse(w, http.StatusNotFound, response) // Corrected to 404 Not Found
			} else {
				data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
					h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
					h.Services.HomeTemplateData.WithUserID(h.Session.GetAuthUserID(r.Context())),
				)

				if templateDataErr != nil {
					h.Loggers.Error.Printf("ERROR: AddCartItemHandler - AddCartItem - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
					http.Error(w, "Failed to load page content", http.StatusInternalServerError)
					return
				}
				data.Errors = append(data.Errors, notFoundErr.Error())
				h.RenderHtmlTemplate(w, "index.html", data, http.StatusNotFound)
			}
		} else {
			// Other internal errors from AddCartItem (e.g., unexpected DB error)
			h.Loggers.Error.Printf("ERROR: Failed to add item to cart (ProductID: %d): %v", addOrder.ProductID, err)
			if strings.HasPrefix(acceptType, "application/json") {
				response := responses.NewErrorJsonResponse("An internal error occurred while adding to cart.")
				responses.WriteJsonResponse(w, http.StatusInternalServerError, response) // 500 Internal Server Error
			} else {
				data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
					h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
					h.Services.HomeTemplateData.WithUserID(userID),
				)

				if templateDataErr != nil {
					h.Loggers.Error.Printf("ERROR: AddCartItemHandler - Other Internal Errors - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
					http.Error(w, "Failed to load page content", http.StatusInternalServerError)
					return
				}
				data.Errors = append(data.Errors, "An internal error occurred while adding to cart.")
				h.RenderHtmlTemplate(w, "index.html", data, http.StatusInternalServerError)
			}
		}
		return
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Added item from Cart")
		responses.WriteJsonResponse(w, http.StatusCreated, response)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *WebHandler) RemoveSingleCartItemHandler(w http.ResponseWriter, r *http.Request) {
	pathId := chi.URLParam(r, "product_id")
	productID, err := strconv.Atoi(pathId)
	acceptType := r.Header.Get("Accept")

	csrfToken := h.Session.GetCsrfToken(r.Context())
	userID := h.Session.GetAuthUserID(r.Context())

	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {

			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: RemoveSingleCartItemHandler - Begin - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusBadRequest)
		}
		return
	}

	err = h.Services.CartItem.RemoveSingleCartItem(userID, productID)
	if err != nil || err == models.ErrCartItemNotFound {
		if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			if strings.HasPrefix(acceptType, "application/json") {
				response := responses.NewErrorJsonResponse(notFoundErr.Error())
				responses.WriteJsonResponse(w, http.StatusNotFound, response)
			} else {
				data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
					h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
					h.Services.HomeTemplateData.WithUserID(userID),
				)

				if templateDataErr != nil {
					h.Loggers.Error.Printf("ERROR: RemoveSingleCartItemHandler - Failed to Remove - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
					http.Error(w, "Failed to load page content", http.StatusInternalServerError)
					return
				}
				data.Errors = append(data.Errors, notFoundErr.Error())
				h.RenderHtmlTemplate(w, "index.html", data, http.StatusNotFound)
			}
			return
		}
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusInternalServerError, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: RemoveSingleCartItemHandler - Failed to Close - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusInternalServerError)
		}
		return
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Product Item Removed")
		responses.WriteJsonResponse(w, http.StatusOK, response)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *WebHandler) RemoveCartItemHandler(w http.ResponseWriter, r *http.Request) {
	pathId := chi.URLParam(r, "item_id")
	cartItemId, err := strconv.Atoi(pathId)

	acceptType := r.Header.Get("Accept")

	csrfToken := h.Session.GetCsrfToken(r.Context())
	userID := h.Session.GetAuthUserID(r.Context())

	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: RemoveCartItemHandler - Begin - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusBadRequest)
		}
		return
	}

	err = h.Services.CartItem.RemoveCartItem(userID, cartItemId)
	if err != nil {
		if notFoundErr, ok := err.(*appErrors.NotFoundError); ok || err == models.ErrCartItemNotFound {
			if strings.HasPrefix(acceptType, "application/json") {
				response := responses.NewErrorJsonResponse(notFoundErr.Error())
				responses.WriteJsonResponse(w, http.StatusNotFound, response)
			} else {
				data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
					h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
					h.Services.HomeTemplateData.WithUserID(userID),
				)

				if templateDataErr != nil {
					h.Loggers.Error.Printf("ERROR: RemoveCartItemHandler - Failed to Remove - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
					http.Error(w, "Failed to load page content", http.StatusInternalServerError)
					return
				}
				data.Errors = append(data.Errors, notFoundErr.Error())
				h.RenderHtmlTemplate(w, "index.html", data, http.StatusNotFound)
			}
			return
		}
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusInternalServerError, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: RemoveCartItemHandler - Failed to Close - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusInternalServerError)
		}
		return
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Cart Item Removed")
		responses.WriteJsonResponse(w, http.StatusOK, response)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *WebHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {

	acceptType := r.Header.Get("Accept")

	csrfToken := h.Session.GetCsrfToken(r.Context())
	userID := h.Session.GetAuthUserID(r.Context())

	err := h.Services.CartItem.Checkout(userID)

	if err != nil {
		statusCode := http.StatusBadRequest
		if err == models.ErrCartItemNotFound {
			statusCode = http.StatusNotFound
		}

		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, statusCode, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: CheckoutHandler - Begin - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, statusCode)
		}
		return
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Checked Out.")
		responses.WriteJsonResponse(w, http.StatusOK, response)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *WebHandler) ConfirmOrderHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")

	csrfToken := h.Session.GetCsrfToken(r.Context())
	userID := h.Session.GetAuthUserID(r.Context())

	cart, err := h.Services.CartItem.GetCart(userID)
	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: ConfirmOrderHandler - Begin - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusBadRequest)
		}
		return
	}

	canConfirm := len(cart) > 0

	if !canConfirm {
		err = errors.New("please add some items into your cart")
	}

	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {
			data, templateDataErr := h.Services.HomeTemplateData.GetHomeTemplateContent(
				h.Services.HomeTemplateData.WithCsrfToken(csrfToken),
				h.Services.HomeTemplateData.WithUserID(userID),
			)

			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: ConfirmOrderHandler - Failed to Checkout - Failed to get HTML template content for user %d: %v", userID, templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "index.html", data, http.StatusBadRequest)
		}
		return
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Order Confirmed")
		responses.WriteJsonResponse(w, http.StatusOK, response)
	} else {
		http.Redirect(w, r, "/#confirm-order", http.StatusSeeOther)
	}
}
