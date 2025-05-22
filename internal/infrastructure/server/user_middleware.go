package server

import (
	"github.com/gin-gonic/gin"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
	"github.com/martiriera/discogs-spotify/internal/usecases"
)

func authUserMiddleware(uc usecases.GetSpotifyUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, exists := GetContextValue(ctx, session.SpotifyUserIDKey); exists {
			ctx.Next()
			return
		}

		userID, err := uc.GetUserID(ctx)

		if err != nil || userID == "" {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		SetContextValue(ctx, session.SpotifyUserIDKey, userID)
		ctx.Next()
	}
}
