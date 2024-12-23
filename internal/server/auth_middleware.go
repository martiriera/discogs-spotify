package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"golang.org/x/oauth2"
)

func AuthMiddleware(store session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := store.GetData(c.Request, session.SpotifyTokenKey)
		if err != nil {
			// TODO: log error
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, ok := data.(*oauth2.Token)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set(session.SpotifyTokenKey, token)
		c.Next()
	}
}
