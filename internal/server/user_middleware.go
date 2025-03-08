package server

import (
	"github.com/gin-gonic/gin"

	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/session"
)

func authUserMiddleware(controller spotify.UserController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, exists := ctx.Get(session.SpotifyUserIDKey); exists {
			ctx.Next()
			return
		}

		userID, err := controller.GetSpotifyUserID(ctx)

		if err != nil || userID == "" {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		ctx.Set(session.SpotifyUserIDKey, userID)
		ctx.Next()
	}
}
