package handlers

import (
	"log"
	"net/http"
	"strings"

	"dessert-ordering-go-system/internal/app"
	appConstants "dessert-ordering-go-system/internal/app_constants"
	responses "dessert-ordering-go-system/internal/response"
)

type WebHandler struct {
	*app.Application
}

func NewWebHandlers(a *app.Application) *WebHandler {
	return &WebHandler{a}
}

func (h *WebHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")
	csrfToken := h.Session.GetString(r.Context(), appConstants.X_CSRF_Token)

	if strings.HasPrefix(acceptType, "application/json") {
		response := responses.NewSuccessJsonResponse("Welcome To The Dessert Ordering System")
		responses.WriteJsonHeadersResponse(w, http.StatusOK, response, map[string]string{appConstants.X_CSRF_Token: csrfToken})
		return
	}
	htmlContent, err := h.Services.HomeTemplateData.GetHomeTemplateContent(h.Services.HomeTemplateData.WithCsrfToken(csrfToken))
	sessionFlashError := h.Session.PopString(r.Context(), appConstants.Flash_Error)
	if sessionFlashError != "" {
		htmlContent.Errors = append(htmlContent.Errors, sessionFlashError)
	}

	if err != nil {
		log.Printf("ERROR: HomeHandler - Failed to get HTML template content for user %d: %v", 1, err)
		http.Error(w, "Failed to load page content", http.StatusInternalServerError)
		return
	}

	h.RenderHtmlTemplate(w, "index.html", htmlContent, http.StatusOK)
}

func (h *WebHandler) RedirectToHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
