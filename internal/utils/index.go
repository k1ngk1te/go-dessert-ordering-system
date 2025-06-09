package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomString returns a URL-safe, base64 encoded
// cryptographically secure random string of a given byte length.
func GenerateRandomString(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func FormatPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}
