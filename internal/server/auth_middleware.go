package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"golang.org/x/oauth2"
)

func authMiddleware(store session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Request can have the token so session may not be needed
		// TODDO: TESTS

		tokenJson, err := store.GetData(ctx.Request, session.SpotifyTokenKey)

		if err != nil {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		var token oauth2.Token
		err = json.Unmarshal([]byte(tokenJson.(string)), &token)

		if err != nil {
			ctx.Redirect(302, "/auth/login")
			ctx.Abort()
			return
		}

		ctx.Set(session.SpotifyTokenKey, token)
		ctx.Next()
	}
}
