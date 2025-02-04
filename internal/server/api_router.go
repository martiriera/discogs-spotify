package server

import (
	"net/http"
	"text/template"

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
	userController     *spotify.UserController
	session            *session.Session
	template           *template.Template
}

func NewApiRouter(
	pc *playlist.PlaylistController,
	uc *spotify.UserController,
	s *session.Session,
	t *template.Template) *ApiRouter {
	router := &ApiRouter{playlistController: pc, userController: uc, session: s, template: t}
	return router
}

func (router *ApiRouter) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/", router.handleMain)
	rg.GET("/home", authTokenMiddleware(*router.session), router.handleMain)
	rg.POST("/playlist",
		authTokenMiddleware(*router.session),
		authUserMiddleware(*router.userController),
		router.handlePlaylistCreate,
	)
	rg.Static("/static", "./static")
}

func (router *ApiRouter) handleMain(ctx *gin.Context) {
	if _, exists := ctx.Get(session.SpotifyTokenKey); exists {
		router.handleHome(ctx)
	} else {
		if err := router.template.ExecuteTemplate(ctx.Writer, "index.html", nil); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError)
		}
	}
}

func (router *ApiRouter) handleHome(ctx *gin.Context) {
	if err := router.template.ExecuteTemplate(ctx.Writer, "home.html", nil); err != nil {
		util.HandleError(ctx, err, http.StatusInternalServerError)
	}
}

func (router *ApiRouter) handlePlaylistCreate(ctx *gin.Context) {
	username := ctx.PostForm("discogs_url")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	pl, err := router.playlistController.CreatePlaylist(ctx, username)
	if err != nil {
		if errors.Cause(err) == discogs.ErrUnauthorized {
			util.HandleError(ctx, err, http.StatusUnauthorized)
			return
		}

		if errors.Cause(err) == playlist.ErrInvalidDiscogsUrl {
			util.HandleError(ctx, err, http.StatusBadRequest)
			return
		}

		if errors.Cause(err) == spotify.ErrUnauthorized {
			ctx.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			return
		}

		util.HandleError(ctx, err, http.StatusInternalServerError)
		return
	}

	responseBody := gin.H{
		"id":               pl.ID,
		"url":              pl.URL,
		"discogs_releases": pl.DiscogsReleases,
		"spotify_albums":   pl.SpotifyAlbums,
	}

	ctx.JSON(http.StatusOK, responseBody)
}
