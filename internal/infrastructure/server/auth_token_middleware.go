package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

func authTokenMiddleware(service ports.SessionPort) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, exists := GetContextValue(ctx, session.SpotifyTokenKey); exists {
			ctx.Next()
			return
		}

		token, err := service.GetData(ctx.Request, session.SpotifyTokenKey)

		if err != nil || token == nil || isExpired(token) {
			ctx.Redirect(http.StatusFound, "/auth/login")
			ctx.Abort()
			return
		}

		SetContextValue(ctx, session.SpotifyTokenKey, token)
		ctx.Next()
	}
}

func isExpired(token any) bool {
	oauth2Token := token.(*oauth2.Token)
	return oauth2Token.Expiry.Before(time.Now())
}
