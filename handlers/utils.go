package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// DecodeJSONBody decodes the JSON request body into the provided 'v' struct.
// It returns an error representing the decoding issue and the appropriate
// HTTP status code to send to the client. If decoding is successful,
// it returns nil for error and http.StatusOK (though the status code
// for success is less critical as the caller will set it).
func JsonBodyDecoder(w http.ResponseWriter, r *http.Request, v any) (int, error) {
	// Limit the size of the request body to prevent abuse
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF):
			// Empty body
			return http.StatusBadRequest, errors.New("request body cannot be empty")
		case errors.As(err, &syntaxError):
			// Malformed JSON syntax
			return http.StatusBadRequest, fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			// Incorrect type for a field
			return http.StatusBadRequest, fmt.Errorf("invalid type for field '%s': expected %s but got %v", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
		case errors.Is(err, io.ErrUnexpectedEOF):
			// Incomplete JSON or connection closed prematurely
			return http.StatusBadRequest, errors.New("request body contains incomplete JSON")
		case err.Error() == "http: request body too large": // Error message from http.MaxBytesReader
			return http.StatusRequestEntityTooLarge, errors.New("request body too large")
		default:
			// Catch-all for other unexpected decoding errors
			// Log this full error on the server side for debugging
			log.Printf("ERROR: Unhandled JSON decoding error: %v", err)
			return http.StatusInternalServerError, errors.New("failed to parse request body")
		}
	}

	// After a successful decode, it's good practice to read the rest of the body
	// to ensure the connection can be reused (HTTP keep-alive).
	// This drains any remaining data in the body.
	defer r.Body.Close() // Ensure the body is closed
	io.Copy(io.Discard, r.Body)

	return http.StatusOK, nil // Or just 'return nil, 0' if you don't want to suggest a status code
}
