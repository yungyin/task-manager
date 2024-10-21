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
