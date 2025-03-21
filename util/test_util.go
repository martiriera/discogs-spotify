package util

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
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
