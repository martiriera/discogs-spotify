package server

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

// GetContextValue retrieves a value from the Gin context using a ContextKey
func GetContextValue(ctx *gin.Context, key session.ContextKey) (interface{}, bool) {
	return ctx.Get(string(key))
}

// SetContextValue sets a value in the Gin context using a ContextKey
func SetContextValue(ctx *gin.Context, key session.ContextKey, value interface{}) {
	ctx.Set(string(key), value)
}

// MustGetContextValue retrieves a value from the Gin context using a ContextKey and panics if not found
func MustGetContextValue(ctx *gin.Context, key session.ContextKey) interface{} {
	return ctx.MustGet(string(key))
}
