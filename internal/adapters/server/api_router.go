package server

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	"github.com/martiriera/discogs-spotify/internal/adapters/session"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/util"
)

type APIRouter struct {
	playlistController *playlist.Controller
	userController     *spotify.UserController
	session            *ports.SessionPort
	template           *template.Template
}

func NewAPIRouter(
	pc *playlist.Controller,
	uc *spotify.UserController,
	s *ports.SessionPort,
	t *template.Template) *APIRouter {
	router := &APIRouter{playlistController: pc, userController: uc, session: s, template: t}
	return router
}

func (router *APIRouter) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/", router.handleMain)
	rg.GET("/home", authTokenMiddleware(*router.session), router.handleMain)
	rg.POST("/playlist",
		authTokenMiddleware(*router.session),
		authUserMiddleware(*router.userController),
		router.handlePlaylistCreate,
	)
	rg.Static("/static", "./static")
}

func (router *APIRouter) handleMain(ctx *gin.Context) {
	if _, exists := ctx.Get(session.SpotifyTokenKey); exists {
		router.handleHome(ctx)
	} else {
		if err := router.template.ExecuteTemplate(ctx.Writer, "index.html", nil); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError)
		}
	}
}

func (router *APIRouter) handleHome(ctx *gin.Context) {
	if err := router.template.ExecuteTemplate(ctx.Writer, "home.html", nil); err != nil {
		util.HandleError(ctx, err, http.StatusInternalServerError)
	}
}

func (router *APIRouter) handlePlaylistCreate(ctx *gin.Context) {
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

		if errors.Cause(err) == playlist.ErrInvalidDiscogsURL {
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
