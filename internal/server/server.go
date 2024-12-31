package server

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

type Server struct {
	*gin.Engine
}

func NewServer(
	playlistController *playlist.PlaylistController,
	oauthController *spotify.OAuthController,
	userController *spotify.UserController,
	session session.Session,
) *Server {
	s := &Server{Engine: gin.Default()}

	apiRouter := NewApiRouter(playlistController, userController, &session)
	authRouter := NewAuthRouter(oauthController, &session)

	authGroup := s.Engine.Group("/auth")
	authRouter.SetupRoutes(authGroup)

	apiGroup := s.Engine.Group("/api")
	apiRouter.SetupRoutes(apiGroup)

	return s
}
