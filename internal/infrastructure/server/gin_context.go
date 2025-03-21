package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

func GetContextValue(ctx *gin.Context, key session.ContextKey) (any, bool) {
	return ctx.Get(string(key))
}

func SetContextValue(ctx *gin.Context, key session.ContextKey, value any) {
	ctx.Set(string(key), value)
}

func MustGetContextValue(ctx *gin.Context, key session.ContextKey) any {
	return ctx.MustGet(string(key))
}

type GinContextProvider struct{}

func NewGinContextProvider() *GinContextProvider {
	return &GinContextProvider{}
}

func (p *GinContextProvider) GetToken(ctx context.Context) (*oauth2.Token, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("context is not a gin.Context")
	}

	value, exists := GetContextValue(ginCtx, session.SpotifyTokenKey)
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	token, ok := value.(*oauth2.Token)
	if !ok {
		return nil, fmt.Errorf("token in context is not of type *oauth2.Token")
	}

	return token, nil
}

func (p *GinContextProvider) GetUserID(ctx context.Context) (string, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return "", fmt.Errorf("context is not a gin.Context")
	}

	value, exists := GetContextValue(ginCtx, session.SpotifyUserIDKey)
	if !exists {
		return "", fmt.Errorf("user ID not found in context")
	}

	userID, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("user ID in context is not of type string")
	}

	return userID, nil
}

func (p *GinContextProvider) SetUserID(ctx context.Context, userID string) error {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return fmt.Errorf("context is not a gin.Context")
	}

	SetContextValue(ginCtx, session.SpotifyUserIDKey, userID)
	return nil
}
