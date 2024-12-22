package server

import (
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

type Server struct {
	http.Handler
}

func NewServer(
	playlistCreator *playlist.PlaylistCreator,
	oauthController *spotify.OAuthController,
) *Server {
	s := new(Server)

	apiRouter := NewApiRouter(playlistCreator)
	authRouter := NewAuthRouter(oauthController)

	combinedRouter := http.NewServeMux()
	combinedRouter.Handle("/api/", http.StripPrefix("/api", apiRouter))
	combinedRouter.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	s.Handler = combinedRouter
	return s
}
