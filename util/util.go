package util

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, err error, statusCode int) {
	log.Println(err)
	errorResponse := map[string]string{
		"error": err.Error(),
	}
	ctx.JSON(statusCode, errorResponse)
}

func AssertEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("environment variable %s is required", name)
	}
	return value
}
