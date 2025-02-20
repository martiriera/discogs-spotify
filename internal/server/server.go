package server

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
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
	playlistController *playlist.PlaylistController,
	oauthController *spotify.OAuthController,
	userController *spotify.UserController,
	session session.Session,
) *Server {
	s := &Server{Engine: gin.Default()}

	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))

	apiRouter := NewApiRouter(playlistController, userController, &session, tmpl)
	authRouter := NewAuthRouter(oauthController, &session)

	authGroup := s.Engine.Group("/auth")
	authRouter.SetupRoutes(authGroup)

	apiGroup := s.Engine.Group("/")
	apiRouter.SetupRoutes(apiGroup)

	return s
}
