package server

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
)

func authMiddleware(store session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Request can have the token so session may not be needed
		// TODDO: TESTS

		token, err := store.GetData(ctx.Request, session.SpotifyTokenKey)

		if err != nil {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		ctx.Set(session.SpotifyTokenKey, token)
		ctx.Next()
	}
}
