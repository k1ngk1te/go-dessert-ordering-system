# 🍰 Dessert Order System – Go & Chi

A simple dessert ordering system built with Go and the `chi` router. It supports basic product listing, shopping cart operations, and checkout functionality, all maintained in memory (no external database).

## 🚀 Features

- List available desserts
- Add items to a shopping cart
- View cart contents
- Remove items from cart
- Checkout and clear cart

# Dessert Ordering System (Go)

[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org/doc/go1.22)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/your-username/your-repo-name/actions) ## Table of Contents

- [Dessert Ordering System (Go)](#dessert-ordering-system-go)
  - [Table of Contents](#table-of-contents)
  - [1. Overview](#1-overview)
  - [2. Features](#2-features)
  - [3. Technologies Used](#3-technologies-used)
  - [4. Architecture & Design Principles](#4-architecture--design-principles)
  - [5. Getting Started](#5-getting-started)
    - [Prerequisites](#prerequisites)
    - [Environment Variables](#environment-variables)
    - [Database Setup (MySQL)](#database-setup-mysql)
    - [Redis Setup](#redis-setup)
    - [Installation](#installation)
    - [Running the Application](#running-the-application)
  - [6. API Endpoints Overview](#6-api-endpoints-overview)
  - [7. Authentication](#7-authentication)
  - [8. Folder Structure](#8-folder-structure)
  - [9. Contributing](#9-contributing)
  - [10. License](#10-license)
  - [11. Contact](#11-contact)

## 1. Overview

This project is a robust and secure web application for a dessert ordering system, built entirely with **Go (Golang)**. It's designed to cater to both traditional web clients (rendering HTML views) and modern Single Page Applications (SPAs) or mobile clients (via a RESTful JSON API). The system provides core functionalities for user management, product Browse, and shopping cart operations, all while prioritizing security and maintainability.

## 2. Features

- **User Authentication:**
  - User Registration and Login (via web forms and JSON API).
  - Secure password hashing (`bcrypt`).
  - **Hybrid Authentication:** Supports both secure server-side sessions (for HTML views using `scs` with Redis) and stateless JWT (JSON Web Token) authentication (for API requests).
  - Secure Logout.
- **Product Management:** Browse available desserts.
- **Shopping Cart:** Add/remove items from a cart.
- **CSRF Protection:** Implemented for all HTML forms to prevent Cross-Site Request Forgery attacks.
- **Flash Messages:** Provides user feedback across redirects (e.g., successful login, error messages).
- **Structured Logging:** Centralized logging for errors and application events.
- **Dependency Injection:** A clear and maintainable way to manage application services and dependencies.

## 3. Technologies Used

- **Go (Golang):** The primary programming language.
- **MySQL:** Relational database for storing user accounts, products, orders, etc.
- **Redis:** In-memory data store used for session management (`scs`) and caching.
- **`scs` (Secure Cookie Sessions):** For robust, tamper-proof session management in traditional web flows.
- **`golang-jwt/jwt/v5`:** For generating and validating JSON Web Tokens for API authentication.
- **`bcrypt`:** For secure one-way password hashing.
- **`DotEnv`:** For managing environment variables.
- **Docker & Docker Compose:** For easy setup and management of database and Redis services.

## 4. Architecture & Design Principles

The application adheres to a modular and layered architecture, promoting separation of concerns:

- **Handlers (Controllers):** Handle incoming HTTP requests, delegate logic to services, and render responses (HTML or JSON).
- **Services (Business Logic):** Contain the core business rules and orchestrate interactions between models.
- **Models (Data Access Objects):** Interact directly with the database, abstracting data storage details.
- **Middleware:** Intercepts requests for cross-cutting concerns like authentication, logging, and CSRF protection.
- **Dependency Injection:** Services and components are passed as dependencies, enhancing testability and flexibility.

## 5. Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- **Go:** Version 1.22 or higher.
- **Git:** For cloning the repository.
- **Docker & Docker Compose:** (Recommended) For easily setting up MySQL and Redis. Alternatively, you can install and configure them manually.

### Environment Variables

Create a `.env` file in the root directory of the project and populate it with the following environment variables:

```dotenv
# Application Port
PORT=8080

# Database Configuration (MySQL)
DB_DSN="user:password@tcp(localhost:3306)/dessert_db?parseTime=true"
# Example if using Docker Compose:
# DB_DSN="user:password@tcp(mysql:3306)/dessert_db?parseTime=true"

# Redis Configuration
REDIS_ADDR="localhost:6379"
# Example if using Docker Compose:
# REDIS_ADDR="redis:6379"
REDIS_PASSWORD="" # Leave empty if no password
REDIS_DB=0

# JWT Secret Key (for JSON Web Tokens)
JWT_SECRET_KEY="your_very_long_and_random_jwt_secret_key"

# Application Domain (e.g., for cookie security)
APP_DOMAIN="localhost"
```
