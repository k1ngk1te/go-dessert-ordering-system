# Dessert Ordering System (GoLang)

A robust and scalable backend system for managing dessert orders, built with Go. This application features a dual-authentication mechanism (JWT for APIs and Sessions for web), comprehensive input validation, and a clear, modular architecture.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Setup and Installation](#setup-and-installation)
- [Running the Application](#running-the-application)
- [API Endpoints (Examples)](#api-endpoints-examples)
- [Authentication Details](#authentication-details)
- [Input Validation](#input-validation)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Mixed Authentication:** Seamlessly supports both JSON Web Token (JWT) based authentication for API clients (via HttpOnly cookies and `Authorization: Bearer` headers) and traditional session-based authentication for web clients.
- **User Management:** Secure user registration, authentication, and session management.
- **Robust Input Validation:** Comprehensive server-side validation for all incoming user data (forms and API requests) using `github.com/go-playground/validator`. Provides detailed error feedback.
- **Structured Error Handling:** Differentiates between JSON error responses for API consumers and HTML redirects with flash messages for web UIs.
- **Context-Based User Data:** Authenticated user information (ID, username) is securely passed through the request context.
- **Logging:** Integrated logging for better observability and debugging.
- **Flash Messages:** User-friendly feedback on web pages for actions like login failures or successful operations.
- **CSRF Protection:** (Crucial when using HttpOnly JWT cookies alongside session management - ensure this is properly implemented in relevant areas, e.g., via `scs`'s built-in CSRF token or a custom mechanism for API forms).

## Technologies Used

- **Go (Golang):** The primary programming language.
- **HTTP Router:** (e.g., `github.com/go-chi/chi` or `gorilla/mux`) for handling routes.
- **JSON Web Tokens (JWT):** `github.com/golang-jwt/jwt/v5` for API authentication.
- **Sessions:** `github.com/alexedwards/scs/v2` (or similar) for session management and flash messages.
- **Input Validation:** `github.com/go-playground/validator/v10` for powerful struct-based validation.
- **Database:** (e.g., PostgreSQL, MySQL) for data persistence.
- **Logging:** Standard library `log` or a structured logging library.

## Project Structure

The project is organized into logical packages reflecting different layers and concerns of the application, promoting modularity and maintainability.

/project-root
├── app/ # Contains the core Application struct for shared resources (Logger, DB, Config, Session)
│ └── application.go
│ └── context_keys.go # Custom type for context keys (e.g., for UserID)
│
├── controllers/ # Houses HTTP handlers / controllers for various application interfaces
│ └── web/ # Web-specific HTTP handlers responsible for interacting with HTML views
│ └── web_handlers.go
│ └── routes.go # (If routes are defined in a separate file within controllers/web)
│
├── services/ # Contains core business logic and service implementations
│ ├── auth_service.go # Handles authentication logic (JWT, session, user authentication)
│ ├── user_service.go # Manages user-related business logic
│ ├── html_template_service.go # Service for rendering HTML templates
│ └── (other_services).go # e.g., product_service.go, order_service.go
│
├── internal/ # A collection of internal utility and helper packages,
│ │ # typically not intended for direct import by other Go modules
│ ├── app_constants/ # Defines application-wide constants
│ │ └── constants.go
│ │
│ ├── response/ # Provides standardized HTTP response structures (e.g., for JSON error responses)
│ │ └── response.go
│ │
│ └── utils/ # Contains general utility functions
│ └── helpers.go
│ └── validation/ # Holds the custom validator setup and related helper functions
│ └── validator.go
│
├── main.go # The main application entry point, responsible for bootstrapping the server
├── go.mod # Go module definition file, managing project dependencies
└── go.sum # Go module checksums file, for verifying module authenticity
