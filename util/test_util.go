package util

import (
	"context"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
	"golang.org/x/oauth2"
)

func NewTestContextWithToken(key session.ContextKey, token *oauth2.Token) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, key, token)
}

// func NewTestGinContextWithToken(key string, token *oauth2.Token) *gin.Context {
// 	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
// 	ctx.Set(key, token)
// 	return ctx
// }
