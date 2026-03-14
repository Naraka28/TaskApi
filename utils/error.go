package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct{
	Error string `json:"error"`
	Message string `json:"message"`
}

func SendJSONError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)

    response := ErrorResponse{
        Error:   http.StatusText(code), // Pone "Not Found", "Internal Server Error", etc.
        Message: message,
    }

    json.NewEncoder(w).Encode(response)
}