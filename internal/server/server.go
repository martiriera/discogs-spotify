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

	// TODO: Implement the OAuthController.SetTokenStoreFunc method
	// oauthController.SetTokenStoreFunc(func(token *oauth2.Token) {
	// 	fmt.Println("Access Token:", token.AccessToken)
	// 	fmt.Println("Refresh Token:", token.RefreshToken)
	// })

	s.Handler = combinedRouter
	return s
}
