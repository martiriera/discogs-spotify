package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
	"github.com/pkg/errors"
)

type ApiRouter struct {
	playlistCreator *playlist.PlaylistCreator
}

func NewApiRouter(c *playlist.PlaylistCreator) *ApiRouter {
	router := &ApiRouter{playlistCreator: c}
	router.playlistCreator = c
	return router
}

func (router *ApiRouter) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/", router.handleMain)
	rg.POST("/playlist", router.handlePlaylistCreate)
}

func (router *ApiRouter) handleMain(c *gin.Context) {
	html := `<html>
					<body>
						<form action="/playlist" method="get">
							<label for="username">Discogs username:</label>
							<input type="text" id="username" name="username">
							<button type="submit">Create playlist</button>
						</form>
					</body>
				</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (router *ApiRouter) handlePlaylistCreate(c *gin.Context) {
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	uris, err := router.playlistCreator.CreatePlaylist(username)
	if err != nil {
		if errors.Cause(err) == discogs.ErrUnauthorized {
			util.HandleError(c, err, http.StatusUnauthorized)
			return
		}

		if errors.Cause(err) == spotify.ErrUnauthorized {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		}

		util.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, uris)
}
