package handlers

import (
	"net/http"
	"strings"

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
	acceptType := w.Header().Get("Accept")
	contentType := w.Header().Get("Content-Type")

	formData := &services.LoginForm{}

	if strings.HasPrefix(contentType, "application/json") {
		errStatusCode, err := JsonBodyDecoder(w, r, formData)
		if err != nil {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, errStatusCode, response)
			return
		}
	} else {
		formData.Contact = r.FormValue("contact")
		formData.Password = r.FormValue("password")
	}

	userData, err := h.Services.Auth.Authenticate(formData.Contact, formData.Password)
	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
		} else {
			csrfToken := h.Session.GetCsrfToken(r.Context())
			data, templateDataErr := h.Services.LoginTemplateData.GetLoginTemplateContent(h.Services.LoginTemplateData.WithCsrfToken(csrfToken))
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostLoginHandler - GetLoginTemplateContent: %v", templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Form = formData
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "login.html", data, http.StatusUnauthorized)
		}
		return
	}

	h.Session.SetAuthUserID(r.Context(), userData.ID)
	csrfToken := h.Session.GetCsrfToken(r.Context())
	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonDataResponse("Log in successful", userData)
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
	acceptType := w.Header().Get("Accept")
	contentType := w.Header().Get("Content-Type")

	formData := &services.RegisterForm{}

	if strings.HasPrefix(contentType, "application/json") {
		errStatusCode, err := JsonBodyDecoder(w, r, formData)
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

	err := h.Services.Auth.RegisterUser(formData.Username, formData.Email, formData.Password)
	if err != nil {
		if strings.HasPrefix(acceptType, "application/json") {
			response := responses.NewErrorJsonResponse(err.Error())
			responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
		} else {
			csrfToken := h.Session.GetCsrfToken(r.Context())
			data, templateDataErr := h.Services.RegisterTemplateData.GetRegisterTemplateContent(h.Services.RegisterTemplateData.WithCsrfToken(csrfToken))
			if templateDataErr != nil {
				h.Loggers.Error.Printf("ERROR: PostRegisterHandler - GetRegisterTemplateContent: %v", templateDataErr)
				http.Error(w, "Failed to load page content", http.StatusInternalServerError)
				return
			}
			data.Form = formData
			data.Errors = append(data.Errors, err.Error())
			h.RenderHtmlTemplate(w, "register.html", data, http.StatusUnauthorized)
		}
		return
	}

	csrfToken := h.Session.GetCsrfToken(r.Context())
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
