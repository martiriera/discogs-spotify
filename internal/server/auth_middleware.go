package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"golang.org/x/oauth2"
)

func authMiddleware(store session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := store.GetData(ctx.Request, session.SpotifyTokenKey)
		if err != nil {
			// TODO: log error
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, ok := data.(*oauth2.Token)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Set(session.SpotifyTokenKey, token)
		ctx.Next()
	}
}
