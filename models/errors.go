package models

import (
	"errors"
	"strings"
)

var (
	// Common
	ErrDuplicateRecord = errors.New("duplicate record found")
	// Product
	ErrProductNotFound = errors.New("product not found")
	// Cart
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrNoCartItemsFound = errors.New("no cart items found")
	// User
	ErrInvalidCredentials = errors.New("invalid authentication credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateUsername  = errors.New("duplicate username")
)

// IsDuplicateEntryError is a helper function to check for duplicate entry errors.
// This example is for MySQL; error codes differ for other databases (e.g., PostgreSQL).
func IsDuplicateEntryError(err error) bool {
	// MySQL's duplicate entry error code is 1062.
	// You might need to import "github.com/go-sql-driver/mysql" if you are using MySQL.
	// And then cast the error to *mysql.MySQLError.
	// For now, a simplified check.
	// A more robust check might involve parsing the error message or specific error types.
	return strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062")
}
