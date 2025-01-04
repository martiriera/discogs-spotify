package util

import (
	"log"
	"os"
	"time"

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

func StartTimer(name string) func() {
	t := time.Now()
	return func() {
		d := time.Since(t)
		log.Println(name, "took", d)
	}
}
