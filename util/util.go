package util

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// HandleError formats and sends the error response using Gin
func HandleError(ctx *gin.Context, err error, statusCode int) {
	// Log the error (you can add more sophisticated logging)
	log.Println(err)

	// Prepare the error response
	errorResponse := map[string]string{
		"error": err.Error(),
	}

	// Set response headers and send the error response
	ctx.JSON(statusCode, errorResponse)
}

func AssertEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("environment variable %s is required", name)
	}
	return value
}
