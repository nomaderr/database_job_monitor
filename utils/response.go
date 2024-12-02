package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}
