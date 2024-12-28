package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"golang.org/x/oauth2"
)

func authMiddleware(store session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODDO: TESTS
		if _, exists := ctx.Get(session.SpotifyTokenKey); exists {
			ctx.Next()
			return
		}

		token, err := store.GetData(ctx.Request, session.SpotifyTokenKey)

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
