package server

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/config"
	"github.com/martiriera/discogs-spotify/internal/usecases"
)

type Server struct {
	*gin.Engine
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Engine.ServeHTTP(w, r)
}

//go:embed templates/*
var templateFS embed.FS

func NewServer(
	playlistController *usecases.Controller,
	authenticateSpotify *usecases.SpotifyAuthenticate,
	getSpotifyUser *usecases.GetSpotifyUser,
	session ports.SessionPort,
	cfg *config.Config,
) *Server {
	s := &Server{Engine: gin.Default()}

	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))

	apiRouter := NewAPIRouter(playlistController, getSpotifyUser, &session, tmpl)
	authRouter := NewAuthRouter(authenticateSpotify, &session, cfg)

	authGroup := s.Engine.Group("/auth")
	authRouter.SetupRoutes(authGroup)

	apiGroup := s.Engine.Group("/")
	apiRouter.SetupRoutes(apiGroup)

	return s
}
