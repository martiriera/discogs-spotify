package server

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	coreErrors "github.com/martiriera/discogs-spotify/internal/core/errors"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
	"github.com/martiriera/discogs-spotify/internal/usecases"
)

type APIRouter struct {
	playlistController *usecases.Controller
	userController     *usecases.GetSpotifyUser
	session            *ports.SessionPort
	template           *template.Template
}

func NewAPIRouter(
	pc *usecases.Controller,
	getSpotifyUserUseCase *usecases.GetSpotifyUser,
	session *ports.SessionPort,
	template *template.Template) *APIRouter {
	router := &APIRouter{playlistController: pc, userController: getSpotifyUserUseCase, session: session, template: template}
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
	if _, exists := GetContextValue(ctx, session.SpotifyTokenKey); exists {
		router.handleHome(ctx)
	} else {
		if err := router.template.ExecuteTemplate(ctx.Writer, "index.html", nil); err != nil {
			handleError(ctx, err, http.StatusInternalServerError)
		}
	}
}

func (router *APIRouter) handleHome(ctx *gin.Context) {
	if err := router.template.ExecuteTemplate(ctx.Writer, "home.html", nil); err != nil {
		handleError(ctx, err, http.StatusInternalServerError)
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
			handleError(ctx, err, http.StatusUnauthorized)
			return
		}

		if errors.Cause(err) == usecases.ErrInvalidDiscogsURL {
			handleError(ctx, err, http.StatusBadRequest)
			return
		}

		if errors.Cause(err) == coreErrors.ErrUnauthorized {
			ctx.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			return
		}

		handleError(ctx, err, http.StatusInternalServerError)
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
