package server

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

type Server struct {
	*gin.Engine
}

func NewServer(
	playlistController *playlist.PlaylistController,
	oauthController *spotify.OAuthController,
) *Server {
	s := &Server{Engine: gin.Default()}

	apiRouter := NewApiRouter(playlistController)
	authRouter := NewAuthRouter(oauthController)

	authGroup := s.Engine.Group("/auth")
	authRouter.SetupRoutes(authGroup)

	apiGroup := s.Engine.Group("/api")
	apiRouter.SetupRoutes(apiGroup)

	return s
}
