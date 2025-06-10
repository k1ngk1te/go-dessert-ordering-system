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

* **Mixed Authentication:** Seamlessly supports both JSON Web Token (JWT) based authentication for API clients (via HttpOnly cookies and `Authorization: Bearer` headers) and traditional session-based authentication for web clients.
* **User Management:** Secure user registration, authentication, and session management.
* **Robust Input Validation:** Comprehensive server-side validation for all incoming user data (forms and API requests) using `github.com/go-playground/validator`. Provides detailed error feedback.
* **Structured Error Handling:** Differentiates between JSON error responses for API consumers and HTML redirects with flash messages for web UIs.
* **Context-Based User Data:** Authenticated user information (ID, username) is securely passed through the request context.
* **Logging:** Integrated logging for better observability and debugging.
* **Flash Messages:** User-friendly feedback on web pages for actions like login failures or successful operations.
* **CSRF Protection:** (Crucial when using HttpOnly JWT cookies alongside session management - ensure this is properly implemented in relevant areas, e.g., via `scs`'s built-in CSRF token or a custom mechanism for API forms).

## Technologies Used

* **Go (Golang):** The primary programming language.
* **HTTP Router:** (e.g., `github.com/go-chi/chi` or `gorilla/mux`) for handling routes.
* **JSON Web Tokens (JWT):** `github.com/golang-jwt/jwt/v5` for API authentication.
* **Sessions:** `github.com/alexedwards/scs/v2` (or similar) for session management and flash messages.
* **Input Validation:** `github.com/go-playground/validator/v10` for powerful struct-based validation.
* **Database:** (e.g., PostgreSQL, MySQL) for data persistence.
* **Logging:** Standard library `log` or a structured logging library.

## Project Structure

The project is organized into logical packages reflecting different layers and concerns of the application, promoting modularity and maintainability.

/project-root
├── app/             # Contains the core Application struct for shared resources (Logger, DB, Config, Session)
│   └── application.go
│   └── context_keys.go # Custom type for context keys (e.g., for UserID)
│
├── controllers/     # Houses HTTP handlers / controllers for various application interfaces
│   └── web/         # Web-specific HTTP handlers responsible for interacting with HTML views
│       └── web_handlers.go
│       └── routes.go # (If routes are defined in a separate file within controllers/web)
│
├── services/        # Contains core business logic and service implementations
│   ├── auth_service.go    # Handles authentication logic (JWT, session, user authentication)
│   ├── user_service.go    # Manages user-related business logic
│   ├── html_template_service.go # Service for rendering HTML templates
│   └── (other_services).go # e.g., product_service.go, order_service.go
│
├── internal/        # A collection of internal utility and helper packages,
│   │                # typically not intended for direct import by other Go modules
│   ├── app_constants/   # Defines application-wide constants
│   │   └── constants.go
│   │
│   ├── response/      # Provides standardized HTTP response structures (e.g., for JSON error responses)
│   │   └── response.go
│   │
│   └── utils/         # Contains general utility functions
│       └── helpers.go
│       └── validation/ # Holds the custom validator setup and related helper functions
│           └── validator.go
│
├── main.go          # The main application entry point, responsible for bootstrapping the server
├── go.mod           # Go module definition file, managing project dependencies
└── go.sum           # Go module checksums file, for verifying module authenticity

## Setup and Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/dessert-ordering-system.git](https://github.com/your-username/dessert-ordering-system.git)
    cd dessert-ordering-system
    ```

2.  **Install Go Modules:**
    ```bash
    go mod tidy
    ```

3.  **Database Setup:**
    * Ensure you have a [Your Database] instance running (e.g., PostgreSQL).
    * Create a database for the project.
    * Run any necessary database migrations (you'll need to implement your migration tool or scripts).

4.  **Environment Variables:**
    Create a `.env` file in the project root or set the following environment variables:
    * `PORT=8080` (or your desired port)
    * `DB_CONNECTION_STRING="host=localhost port=5432 user=youruser password=yourpass dbname=yourdb sslmode=disable"`
    * `JWT_SECRET="supersecretjwtkey"` (Choose a strong, random key)
    * `SESSION_SECRET="anothersupersecretkey"` (Choose a strong, random key for session encryption)

## Running the Application

To run the main API server:

```bash
go run .

Understood! My apologies for overriding your preferred project structure. I'll regenerate the README.md and ensure the "Project Structure" section reflects the layout you had in mind, based on the package paths and components we've discussed throughout our conversation.

Here's the regenerated README.md with the "Project Structure" section left as per your instruction, focusing on how your existing modules are organized:

Markdown

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

* **Mixed Authentication:** Seamlessly supports both JSON Web Token (JWT) based authentication for API clients (via HttpOnly cookies and `Authorization: Bearer` headers) and traditional session-based authentication for web clients.
* **User Management:** Secure user registration, authentication, and session management.
* **Robust Input Validation:** Comprehensive server-side validation for all incoming user data (forms and API requests) using `github.com/go-playground/validator`. Provides detailed error feedback.
* **Structured Error Handling:** Differentiates between JSON error responses for API consumers and HTML redirects with flash messages for web UIs.
* **Context-Based User Data:** Authenticated user information (ID, username) is securely passed through the request context.
* **Logging:** Integrated logging for better observability and debugging.
* **Flash Messages:** User-friendly feedback on web pages for actions like login failures or successful operations.
* **CSRF Protection:** (Crucial when using HttpOnly JWT cookies alongside session management - ensure this is properly implemented in relevant areas, e.g., via `scs`'s built-in CSRF token or a custom mechanism for API forms).

## Technologies Used

* **Go (Golang):** The primary programming language.
* **HTTP Router:** (e.g., `github.com/go-chi/chi` or `gorilla/mux`) for handling routes.
* **JSON Web Tokens (JWT):** `github.com/golang-jwt/jwt/v5` for API authentication.
* **Sessions:** `github.com/alexedwards/scs/v2` (or similar) for session management and flash messages.
* **Input Validation:** `github.com/go-playground/validator/v10` for powerful struct-based validation.
* **Database:** (e.g., PostgreSQL, MySQL) for data persistence.
* **Logging:** Standard library `log` or a structured logging library.

## Project Structure

The project is organized into logical packages reflecting different layers and concerns of the application, promoting modularity and maintainability.

/project-root
├── app/             # Contains the core Application struct for shared resources (Logger, DB, Config, Session)
│   └── application.go
│   └── context_keys.go # Custom type for context keys (e.g., for UserID)
│
├── controllers/     # Houses HTTP handlers / controllers for various application interfaces
│   └── web/         # Web-specific HTTP handlers responsible for interacting with HTML views
│       └── web_handlers.go
│       └── routes.go # (If routes are defined in a separate file within controllers/web)
│
├── services/        # Contains core business logic and service implementations
│   ├── auth_service.go    # Handles authentication logic (JWT, session, user authentication)
│   ├── user_service.go    # Manages user-related business logic
│   ├── html_template_service.go # Service for rendering HTML templates
│   └── (other_services).go # e.g., product_service.go, order_service.go
│
├── internal/        # A collection of internal utility and helper packages,
│   │                # typically not intended for direct import by other Go modules
│   ├── app_constants/   # Defines application-wide constants
│   │   └── constants.go
│   │
│   ├── response/      # Provides standardized HTTP response structures (e.g., for JSON error responses)
│   │   └── response.go
│   │
│   └── utils/         # Contains general utility functions
│       └── helpers.go
│       └── validation/ # Holds the custom validator setup and related helper functions
│           └── validator.go
│
├── main.go          # The main application entry point, responsible for bootstrapping the server
├── go.mod           # Go module definition file, managing project dependencies
└── go.sum           # Go module checksums file, for verifying module authenticity


## Setup and Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/dessert-ordering-system.git](https://github.com/your-username/dessert-ordering-system.git)
    cd dessert-ordering-system
    ```

2.  **Install Go Modules:**
    ```bash
    go mod tidy
    ```

3.  **Database Setup:**
    * Ensure you have a [Your Database] instance running (e.g., PostgreSQL).
    * Create a database for the project.
    * Run any necessary database migrations (you'll need to implement your migration tool or scripts).

4.  **Environment Variables:**
    Create a `.env` file in the project root or set the following environment variables:
    * `PORT=8080` (or your desired port)
    * `DB_CONNECTION_STRING="host=localhost port=5432 user=youruser password=yourpass dbname=yourdb sslmode=disable"`
    * `JWT_SECRET="supersecretjwtkey"` (Choose a strong, random key)
    * `SESSION_SECRET="anothersupersecretkey"` (Choose a strong, random key for session encryption)

## Running the Application

To run the main API server:

```bash
go run .
The server should start on the port specified in your PORT environment variable (defaulting to 8080 if not set).

API Endpoints (Examples)
(You would list your main API endpoints here. For example:)

POST /api/v1/register: User registration.
POST /api/v1/login: User login (returns JWT or sets session cookie).
GET /api/v1/profile: Get authenticated user's profile (requires authentication).
GET /api/v1/desserts: List available desserts.
POST /api/v1/orders: Create a new order (requires authentication).
Authentication Details
The system employs a mixed authentication strategy:

API Clients: Primarily use JWTs. Upon successful login, a JWT is issued and can be stored in an HttpOnly cookie (for browser-based SPAs) or sent in the Authorization: Bearer header (for mobile apps, other services).
Web Clients: Leverage traditional server-side sessions, typically identified by a session cookie.
Hybrid Flow: If a user logs in via JWT, their session is also implicitly established/updated, allowing for seamless transitions between API and web-based interactions within the same application.
Input Validation
All incoming data from HTTP requests (JSON bodies and form data) is rigorously validated server-side using github.com/go-playground/validator.

Validation rules are defined using struct tags (e.g., validate:"required,email,min=8"). When validation fails:

For API requests (Accept: application/json): A 400 Bad Request status is returned with a structured JSON response detailing field-specific errors.
For Web forms: The user is typically redirected back to the form with flash messages containing the validation errors, and the form fields are often repopulated with their previous (invalid) input for convenience.
This ensures data integrity, enhances security, and provides clear feedback to the client.

Contributing
Contributions are welcome! Please fork the repository and open a pull request with your changes.

License
This project is licensed under the MIT License - see the LICENSE file for details.