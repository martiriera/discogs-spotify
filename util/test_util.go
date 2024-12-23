package util

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func NewTestContextWithToken(key string, token *oauth2.Token) *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set(key, token)
	return ctx
}
