package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Server struct {
	playlistCreator *playlist.PlaylistCreator
	oauthController *spotify.OAuthController
	http.Handler
}

func NewServer(
	playlistCreator *playlist.PlaylistCreator,
	oauthController *spotify.OAuthController,
) *Server {
	s := new(Server)

	s.playlistCreator = playlistCreator
	router := http.NewServeMux()

	router.Handle("/create-playlist", http.HandlerFunc(s.handlePlaylistCreate))
	router.Handle("/", http.HandlerFunc(s.handleMain))
	router.Handle("/login", http.HandlerFunc(oauthController.HandleLogin))
	router.Handle("/callback", http.HandlerFunc(oauthController.HandleCallback))

	// TODO: Implement the OAuthController.SetTokenStoreFunc method
	oauthController.SetTokenStoreFunc(func(token *oauth2.Token) {
		fmt.Println("Access Token:", token.AccessToken)
		fmt.Println("Refresh Token:", token.RefreshToken)
	})

	s.Handler = router
	return s
}

func (o *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	html := `<html>
				<body>
					<a href="/login">Log in with Spotify</a>
				</body>
			</html>`
	fmt.Fprint(w, html)
}

func (s *Server) handlePlaylistCreate(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	uris, err := s.playlistCreator.CreatePlaylist(username)
	if err != nil {
		if errors.Cause(err) == discogs.ErrUnauthorized {
			util.HandleError(w, err, http.StatusUnauthorized)
			return
		}

		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(uris)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(response)
}
