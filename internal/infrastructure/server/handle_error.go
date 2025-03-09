package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func handleError(ctx *gin.Context, err error, statusCode int) {
	log.Println(err)
	errorResponse := map[string]string{
		"error": err.Error(),
	}
	ctx.JSON(statusCode, errorResponse)
}
