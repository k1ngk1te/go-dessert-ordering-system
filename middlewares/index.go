package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"slices" // slices.Contains is fine, but direct comparison is also common
	"strings"
	"time"

	app "dessert-ordering-go-system/internal/app"
	appConstants "dessert-ordering-go-system/internal/app_constants"
	responses "dessert-ordering-go-system/internal/response"
	utils "dessert-ordering-go-system/internal/utils"
	services "dessert-ordering-go-system/services"

	"github.com/golang-jwt/jwt/v5"
)

type Middlewares struct {
	*app.Application
}

func NewMiddlewares(a *app.Application) *Middlewares {
	return &Middlewares{a}
}

func (m *Middlewares) EnableCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedRequestMethods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}
		isAllowed := slices.Contains(allowedRequestMethods, r.Method)

		if isAllowed {
			// 1. Get the Token from the session
			sessionID := m.Session.Token(r.Context())
			currentMilliTime := time.Now().UnixMilli()
			if sessionID == "" {
				// 2. It's not available. Add in a new field to initialize it and generate csrf_token
				csrfToken, err := utils.GenerateRandomString(32)
				if err != nil {
					m.Loggers.Error.Printf("failed to generate CSRF Token :%v", err)
				} else {
					m.Session.Put(r.Context(), appConstants.X_CSRF_Token, csrfToken)
				}

				m.Session.Put(r.Context(), "created_at", currentMilliTime)
				m.Session.Put(r.Context(), "last_seen", currentMilliTime)
			} else {
				// Check if the CSRF_TOKEN exists
				csrfToken := m.Session.GetString(r.Context(), appConstants.X_CSRF_Token)
				if csrfToken == "" {
					newCsrfToken, err := utils.GenerateRandomString(32)
					if err != nil {
						m.Loggers.Error.Printf("failed to generate CSRF Token :%v", err)
					} else {
						m.Session.Put(r.Context(), appConstants.X_CSRF_Token, newCsrfToken)
					}
				}
				m.Session.Put(r.Context(), "last_seen", currentMilliTime)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) RequireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedRequestMethods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}
		isAllowed := slices.Contains(allowedRequestMethods, r.Method)

		if !isAllowed { // This block runs for POST, PUT, DELETE, etc.
			// 1. Get the token from the incoming request
			tokenFromRequest := r.FormValue("csrf_token") // From form field
			if tokenFromRequest == "" {
				tokenFromRequest = r.Header.Get(appConstants.X_CSRF_Token) // From custom header
			}

			// 2. Get the expected token from the session
			expectedToken := m.Session.GetCsrfToken(r.Context())

			// 3. Compare the two tokens and check if expected token exists
			if expectedToken == "" || tokenFromRequest == "" || tokenFromRequest != expectedToken {
				m.Loggers.Error.Printf("CSRF token validation failed for request %s %s from %s (Expected: %q, Received: %q)",
					r.Method, r.URL.Path, r.RemoteAddr, expectedToken, tokenFromRequest)

				// 4. Handle forbidden response based on Accept header
				acceptType := r.Header.Get("Accept")
				if strings.HasPrefix(acceptType, "application/json") {
					data := responses.NewErrorJsonResponse("CSRF Token is invalid")
					responses.WriteJsonResponse(w, http.StatusForbidden, data)
				} else {
					// Redirect back to original page with error message, or dedicated error page
					// Ensure your template rendering function handles errors gracefully
					m.Session.Put(r.Context(), appConstants.Flash_Error, "Invalid form submission. Please try again.")
					http.Redirect(w, r, r.URL.Path, http.StatusSeeOther) // Redirect back to same page
					return
				}
				return // Stop execution here if validation fails
			}

			// If tokens match, invalidate the token to prevent replay attacks
			// This is a common and good security practice.
			m.Session.Remove(r.Context(), appConstants.X_CSRF_Token) // Remove the old token
			// Optionally generate a new one immediately for subsequent requests if needed,
			// but usually, it's generated by EnableCSRF on the next GET request.
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptType := r.Header.Get("Accept")

		// JWT Authentication
		var tokenString string = ""

		// 1. Try to get the token from the HttpOnly cookie first
		cookie, err := r.Cookie(appConstants.Jwt_Name)
		if err == nil {
			tokenString = cookie.Value
		} else if errors.Is(err, http.ErrNoCookie) {
			// 2. If not in the cookie, try the Authorization Header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				headerParts := strings.Split(authHeader, " ")
				if len(headerParts) == 2 && headerParts[0] == "Bearer" {
					tokenString = headerParts[1]
				}
			}
		}

		if tokenString != "" {
			claims := &services.UserClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return m.Services.Auth.JWTSecret, nil
			})
			if err != nil {
				if errors.Is(err, jwt.ErrSignatureInvalid) {
					if strings.HasPrefix(acceptType, "application/json") {
						response := responses.NewErrorJsonResponse("Invalid token signature")
						responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
						return
					}

					m.Session.SetFlashError(r.Context(), "Invalid token signature")
					http.Redirect(w, r, "/login", http.StatusSeeOther)

				} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
					if strings.HasPrefix(acceptType, "application/json") {
						response := responses.NewErrorJsonResponse("Token expired or not valid yet")
						responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
						return
					}

					m.Session.SetFlashError(r.Context(), "Token expired or not valid yet")
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				} else {
					if strings.HasPrefix(acceptType, "application/json") {
						response := responses.NewErrorJsonResponse(fmt.Sprintf("Invalid token: %v", err))
						responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
						return
					}
					m.Session.SetFlashError(r.Context(), fmt.Sprintf("Invalid token: %v", err))
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				}
				return
			}

			if !token.Valid {
				if strings.HasPrefix(acceptType, "application/json") {
					response := responses.NewErrorJsonResponse("Invalid token (general validation failure)")
					responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
					return
				}
				m.Session.SetFlashError(r.Context(), "Invalid token (general validation failure)")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Check the CSRF Token
			csrfTokenFromHeader := r.Header.Get(appConstants.X_CSRF_Token)
			expectedCsrfToken := m.Session.GetCsrfToken(r.Context())

			if expectedCsrfToken != csrfTokenFromHeader {
				if strings.HasPrefix(acceptType, "application/json") {
					response := responses.NewErrorJsonResponse("Invalid CSRF Token")
					responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
					return
				}
				m.Session.SetFlashError(r.Context(), "Invalid CSRF Token")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			m.Session.SetAuthUserID(r.Context(), claims.ID)
			next.ServeHTTP(w, r)
			return
		}

		// Session Authentication
		userIDExists := m.Session.Exists(r.Context(), appConstants.Auth_User_ID)
		if !userIDExists {
			if strings.HasPrefix(acceptType, "application/json") {
				response := responses.NewErrorJsonResponse("authentication credentials were not found")
				responses.WriteJsonResponse(w, http.StatusUnauthorized, response)
				return
			}

			m.Session.SetFlashError(r.Context(), "authentication credentials were not found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) AuthNotRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the user ID from the session
		userIDExists := m.Session.Exists(r.Context(), appConstants.Auth_User_ID)
		if userIDExists {
			acceptType := r.Header.Get("Accept")
			if strings.HasPrefix(acceptType, "application/json") {
				response := responses.NewErrorJsonResponse("authentication credentials have been verified")
				responses.WriteJsonResponse(w, http.StatusForbidden, response)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	})
}
