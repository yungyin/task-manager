package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func BadRequestErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Bad Request Error: %v", r)

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusBadRequest)

	response := ErrorResponse{
		Error:   "bad_request_error",
		Message: "The request could not be read properly",
	}

	json.NewEncoder(w).Encode(response)
}

func NotFoundErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Found Error: %v", r)

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusNotFound)

	response := ErrorResponse{
		Error:   "not_found_error",
		Message: "The requested resource was not found.",
	}

	json.NewEncoder(w).Encode(response)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error: %v", err)

	w.Header().Set(ContentType, JsonContentType)
	w.WriteHeader(http.StatusInternalServerError)

	response := ErrorResponse{
		Error:   "internal_server_error",
		Message: "An unexpected error occurred. Please try again later.",
	}

	json.NewEncoder(w).Encode(response)
}
