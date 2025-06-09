package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"dessert-ordering-go-system/handlers"
	"dessert-ordering-go-system/internal/app"
	responses "dessert-ordering-go-system/internal/response"
	middlewares "dessert-ordering-go-system/middlewares"
)

func NewRoutes(a *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// --- Chi's Built-in Middlewares (Commonly used) ---
	r.Use(middleware.Logger)          // Log Requests
	r.Use(middleware.RedirectSlashes) // Support Trailing Slash Requests
	r.Use(middleware.RequestID)       // Adds a request ID to the context
	r.Use(middleware.RealIP)          // Safely extracts the client IP address
	r.Use(middleware.Recoverer)       // Recovers from panics and logs them, prevent server crash
	r.Use(a.Session.LoadAndSave)      // Helps load and save the session automatically

	// --- YOUR CUSTOM 404 HANDLER ---
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// You can log the fact that a 404 occurred for debugging
		a.Loggers.Error.Printf("404 Not Found: %s %s", r.Method, r.URL.Path)

		// Use your existing JSON response structure
		response := responses.NewErrorJsonResponse(fmt.Sprintf("The requested resource '%s' was not found.", r.URL.Path))
		responses.WriteJsonResponse(w, http.StatusNotFound, response) // Use your helper function
	})

	// --- YOUR CUSTOM 405 HANDLER (Optional, but good practice) ---
	// This handles cases where the path exists but the HTTP method is not allowed.
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		a.Loggers.Error.Printf("405 Method Not Allowed: %s %s", r.Method, r.URL.Path)
		response := responses.NewErrorJsonResponse(fmt.Sprintf("Method '%s' not allowed for this resource.", r.Method))
		responses.WriteJsonResponse(w, http.StatusMethodNotAllowed, response)
	})

	// Serve Static Files
	staticDir := http.Dir("./static")

	// Create the file server handler
	fileServer := http.FileServer(staticDir)

	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Initialize custom middlewares
	customMiddlewares := middlewares.NewMiddlewares(a)
	handlers := handlers.NewWebHandlers(a)

	// Authentication Not Required
	r.Group(func(r chi.Router) {
		r.Use(customMiddlewares.AuthNotRequired)

		r.Get("/login", handlers.GetLoginHandler)
		r.Post("/login", handlers.PostLoginHandler)

		r.Get("/register", handlers.GetRegisterHandler)
		r.Post("/register", handlers.PostRegisterHandler)
	})

	// Authentication Required
	r.Group(func(r chi.Router) {
		r.Use(customMiddlewares.AuthRequired)
		r.Group(func(r chi.Router) {
			r.Use(customMiddlewares.EnableCSRF)

			r.Get("/", handlers.HomeHandler)

			r.Get("/products", handlers.GetProductsHandler)
			r.Get("/products/{id}", handlers.GetProductDetailHandler)

			r.Get("/cart", handlers.GetCartHandler)
			r.Get("/cart/product/{product_id}/remove-one", handlers.RedirectToHomeHandler) // Just in case the user refreshes
			r.Get("/cart/{item_id}/delete", handlers.RedirectToHomeHandler)                // Just in case the user refreshes
			r.Get("/confirm-order", handlers.ConfirmOrderHandler)
			r.Get("/checkout", handlers.RedirectToHomeHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(customMiddlewares.RequireCSRF) // Apply require CSRF middlewares to all routes in this group

			r.Get("/logout", handlers.RedirectToHomeHandler)
			r.Post("/logout", handlers.LogoutHandler)
			r.Post("/cart/{item_id}/delete", handlers.RemoveCartItemHandler)
			r.Post("/cart/product/{product_id}/remove-one", handlers.RemoveSingleCartItemHandler)
			r.Post("/cart", handlers.AddCartItemHandler)
			r.Post("/checkout", handlers.CheckoutHandler)
		})
	})

	return r
}
