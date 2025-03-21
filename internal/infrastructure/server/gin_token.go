package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

type GinTokenProvider struct{}

func NewGinTokenProvider() *GinTokenProvider {
	return &GinTokenProvider{}
}

func (p *GinTokenProvider) GetToken(ctx context.Context) (*oauth2.Token, error) {
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
