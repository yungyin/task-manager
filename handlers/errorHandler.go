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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ErrorResponse{
		Error:   "bad_request_error",
		Message: "The request could not be read properly",
	}

	e := json.NewEncoder(w).Encode(response)
	if e != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := ErrorResponse{
		Error:   "internal_server_error",
		Message: "An unexpected error occurred. Please try again later.",
	}

	json.NewEncoder(w).Encode(response)
}
