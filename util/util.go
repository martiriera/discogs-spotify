package util

import (
	"encoding/json"
	"net/http"
)

// HandleError formats and sends the error response
func HandleError(w http.ResponseWriter, err error, statusCode int) {

	// Log the error (you can add more sophisticated logging)
	// Log.Println(wrappedErr)

	// Prepare the error response
	errorResponse := map[string]string{
		"error": err.Error(),
	}

	// Set response headers and send the error response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
