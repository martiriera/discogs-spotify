package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/session"
)

func authTokenMiddleware(service session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, exists := ctx.Get(session.SpotifyTokenKey); exists {
			ctx.Next()
			return
		}

		token, err := service.GetData(ctx.Request, session.SpotifyTokenKey)

		if err != nil || token == nil || isExpired(token) {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		ctx.Set(session.SpotifyTokenKey, token)
		ctx.Next()
	}
}

func isExpired(token any) bool {
	oauth2Token := token.(*oauth2.Token)
	return oauth2Token.Expiry.Before(time.Now())
}
