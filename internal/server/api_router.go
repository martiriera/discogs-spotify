package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
	"github.com/pkg/errors"
)

type ApiRouter struct {
	playlistController *playlist.PlaylistController
	session            *session.Session
}

func NewApiRouter(c *playlist.PlaylistController, s *session.Session) *ApiRouter {
	router := &ApiRouter{playlistController: c, session: s}
	return router
}

func (router *ApiRouter) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/", router.handleMain)
	rg.GET("/home", authMiddleware(*router.session), router.handleMain)
	rg.POST("/playlist", authMiddleware(*router.session), router.handlePlaylistCreate)
}

func (router *ApiRouter) handleMain(ctx *gin.Context) {
	if _, exists := ctx.Get(session.SpotifyTokenKey); exists {
		router.handleHome(ctx)
	} else {
		html := `<html>
					<body>
						<a href="/auth/login">Login with Spotify</a>
					</body>
				</html>`
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

func (router *ApiRouter) handleHome(ctx *gin.Context) {
	html := `<html>
				<body>
					<form action="/api/playlist" method="post">
						<label for="username">Discogs username:</label>
						<input type="text" name="username" name="username">
						<input type="submit" value="Create playlist">
					</form>
				</body>
			</html>`
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (router *ApiRouter) handlePlaylistCreate(ctx *gin.Context) {
	username := ctx.PostForm("username")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	playlistId, err := router.playlistController.CreatePlaylist(ctx, username)
	if err != nil {
		if errors.Cause(err) == discogs.ErrUnauthorized {
			util.HandleError(ctx, err, http.StatusUnauthorized)
			return
		}

		if errors.Cause(err) == spotify.ErrUnauthorized {
			ctx.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			return
		}

		util.HandleError(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"playlist_url": playlistId})
}
