package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JsonDataResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func NewJsonResponse(status string, message string) JsonResponse {
	return JsonResponse{
		Status:  status,
		Message: message,
	}
}

func NewSuccessJsonResponse(message string) JsonResponse {
	return JsonResponse{
		Status:  "success",
		Message: message,
	}
}

func NewErrorJsonResponse(message string) JsonResponse {
	return JsonResponse{
		Status:  "error",
		Message: message,
	}
}

func NewJsonDataResponse[T any](status string, message string, data T) JsonDataResponse[T] {
	return JsonDataResponse[T]{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewSuccessJsonDataResponse[T any](message string, data T) JsonDataResponse[T] {
	return JsonDataResponse[T]{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewErrorJsonDataResponse[T any](message string, data T) JsonDataResponse[T] {
	return JsonDataResponse[T]{
		Status:  "error",
		Message: message,
		Data:    data,
	}
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error: failed to encode JSON response: %v", err)
	}
}

func WriteJsonHeadersResponse(w http.ResponseWriter, statusCode int, data interface{}, headers map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	for name, value := range headers {
		w.Header().Set(name, value)
	}
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error: failed to encode JSON response: %v", err)
	}
}
