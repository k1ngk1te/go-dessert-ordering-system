package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	appConstants "dessert-ordering-go-system/internal/app_constants"
	responses "dessert-ordering-go-system/internal/response"
	utils "dessert-ordering-go-system/internal/utils"
	services "dessert-ordering-go-system/services"
)

// ****** Login Handlers *******

func (h *WebHandler) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := h.Session.GetCsrfToken(r.Context())

	data, templateDataErr := h.Services.LoginTemplateData.GetLoginTemplateContent(h.Services.LoginTemplateData.WithCsrfToken(csrfToken))
	if templateDataErr != nil {
		h.Loggers.Error.Printf("ERROR: GetLoginHandler - GetLoginTemplateContent: %v", templateDataErr)
		http.Error(w, "Failed to load page content", http.StatusInternalServerError)
		return
	}
	flashError := h.Session.PopFlashError(r.Context())
	if flashError != "" {
		data.Errors = append(data.Errors, flashError)
	}
	h.RenderHtmlTemplate(w, "login.html", data, http.StatusOK)
}

func (h *WebHandler) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")

	csrfToken := h.Session.GetCsrfToken(r.Context())

	var formData services.LoginForm

	if r.Header.Get("Content-Type") == "application/json" {
		errStatusCode, err := JsonBodyDecoder(w, r, &formData)
		if err != nil {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, errStatusCode, response)
			return
		}
	} else {
		formData.Contact = r.FormValue("contact")
		formData.Password = r.FormValue("password")
	}

	// -- Perform Validation --
	validationErrors := h.Validator.ValidateStruct(formData)
	if validationErrors != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonDataResponse("Validation failed", validationErrors)
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {
			data, templateDataErr := h.Services.LoginTemplateData.GetLoginTemplateContent(h.Services.LoginTemplateData.WithCsrfToken(csrfToken))
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostLoginHandler - GetLoginTemplateContent on validation error: %v", templateDataErr)
				http.Error(w, "Failed to reload login page after validation error", http.StatusInternalServerError)
				return
			}
			data.Form = &formData

			for field, msg := range validationErrors {
				data.Errors = append(data.Errors, fmt.Sprintf("%s: %s", field, msg))
			}
			h.RenderHtmlTemplate(w, "login.html", data, http.StatusBadRequest) // 400 Bad Request
		}
		return
	}

	userData, err := h.Services.Auth.Authenticate(formData.Contact, formData.Password)
	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
		} else {
			data, templateDataErr := h.Services.LoginTemplateData.GetLoginTemplateContent(h.Services.LoginTemplateData.WithCsrfToken(csrfToken))
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostLoginHandler - GetLoginTemplateContent: %v", templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Form = &formData
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "login.html", data, http.StatusUnauthorized)
		}
		return
	}

	// Log in the user
	h.Session.SetAuthUserID(r.Context(), userData.ID) // Session Auth

	token, err := h.Services.Auth.GenerateAuthToken(userData.ID, userData.Username, userData.Email)
	if err != nil {
		h.Loggers.Error.Printf("ERROR: PostLoginHandler - h.Services.Auth.GenerateAuthToken: %v", err)
		h.Session.SetFlashError(r.Context(), err.Error())
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusInternalServerError, response)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	authData := h.Services.Auth.CreateAuthData(*userData, token)

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonDataResponse("Log in successful", authData)
		secureCookies, _ := appConstants.GetSecureCookies()
		cookie := &http.Cookie{
			Name:     appConstants.Jwt_Name,
			Value:    authData.Token,
			Expires:  time.Now().Add(h.Services.Auth.GetTokenExpiration()),
			HttpOnly: true,
			Secure:   secureCookies,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		responses.WriteJsonHeadersResponse(w, http.StatusOK, response, map[string]string{appConstants.X_CSRF_Token: csrfToken})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ****** Register Handlers *******

func (h *WebHandler) GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := h.Session.GetCsrfToken(r.Context())

	data, templateDataErr := h.Services.RegisterTemplateData.GetRegisterTemplateContent(h.Services.RegisterTemplateData.WithCsrfToken(csrfToken))
	if templateDataErr != nil {
		h.Loggers.Error.Printf("ERROR: GetRegisterHandler - GetRegisterTemplateContent: %v", templateDataErr)
		http.Error(w, "Failed to load page content", http.StatusInternalServerError)
		return
	}
	flashError := h.Session.PopFlashError(r.Context())
	if flashError != "" {
		data.Errors = append(data.Errors, flashError)
	}
	h.RenderHtmlTemplate(w, "register.html", data, http.StatusOK)
}

func (h *WebHandler) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")
	contentType := r.Header.Get("Content-Type")

	csrfToken := h.Session.GetCsrfToken(r.Context())

	var formData services.RegisterForm

	if strings.HasPrefix(contentType, "application/json") {
		errStatusCode, err := JsonBodyDecoder(w, r, &formData)
		if err != nil {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, errStatusCode, response)
			return
		}
	} else {
		formData.Email = r.FormValue("email")
		formData.Username = r.FormValue("username")
		formData.Password = r.FormValue("password")
	}

	// -- Perform Validation --
	validationErrors := h.Validator.ValidateStruct(formData)
	if validationErrors != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonDataResponse("Validation failed", validationErrors)
			responses.WriteJsonResponse(w, http.StatusBadRequest, response)
		} else {
			data, templateDataErr := h.Services.RegisterTemplateData.GetRegisterTemplateContent(h.Services.RegisterTemplateData.WithCsrfToken(csrfToken))
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostRegisterHandler - GetRegisterTemplateContent on validation error: %v", templateDataErr)
				http.Error(w, "Failed to reload register page after validation error", http.StatusInternalServerError)
				return
			}
			data.Form = &formData

			for field, msg := range validationErrors {
				data.Errors = append(data.Errors, fmt.Sprintf("%s: %s", field, msg))
			}
			h.RenderHtmlTemplate(w, "register.html", data, http.StatusBadRequest) // 400 Bad Request
		}
		return
	}

	err := h.Services.Auth.RegisterUser(formData.Username, formData.Email, formData.Password)
	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
		} else {
			data, templateDataErr := h.Services.RegisterTemplateData.GetRegisterTemplateContent(h.Services.RegisterTemplateData.WithCsrfToken(csrfToken))
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostRegisterHandler - GetRegisterTemplateContent: %v", templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Form = &formData
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "register.html", data, http.StatusUnauthorized)
		}
		return
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Registration successful")
		responses.WriteJsonHeadersResponse(w, http.StatusOK, response, map[string]string{appConstants.X_CSRF_Token: csrfToken})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ****** Logout Handlers *******
func (h *WebHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")

	h.Session.Destroy(r.Context())

	newCsrfToken, err := utils.GenerateRandomString(32)
	if err != nil {
		h.Loggers.Error.Printf("ERROR: WebHandler.LogoutHandler - utils.GenerateRandomString: %v", err)
	} else {
		h.Session.SetCsrfToken(r.Context(), newCsrfToken)
	}

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Logout successfully")
		responses.WriteJsonResponse(w, http.StatusOK, response)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
