package server

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

func authUserMiddleware(controller spotify.UserController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, exists := ctx.Get(session.SpotifyUserIdKey); exists {
			ctx.Next()
			return
		}

		userId, err := controller.GetSpotifyUserId(ctx)

		if err != nil || userId == "" {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		ctx.Set(session.SpotifyUserIdKey, userId)
		ctx.Next()
	}
}
