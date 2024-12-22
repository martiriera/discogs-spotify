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
	s.oauthController = oauthController
	router := http.NewServeMux()

	router.Handle("/create-playlist", http.HandlerFunc(s.handlePlaylistCreate))
	router.Handle("/", http.HandlerFunc(s.handleMain))
	router.Handle("/login", http.HandlerFunc(s.handleLogin))
	router.Handle("/callback", http.HandlerFunc(s.handleLoginCallback))

	// TODO: Implement the OAuthController.SetTokenStoreFunc method
	// oauthController.SetTokenStoreFunc(func(token *oauth2.Token) {
	// 	fmt.Println("Access Token:", token.AccessToken)
	// 	fmt.Println("Refresh Token:", token.RefreshToken)
	// })

	s.Handler = router
	return s
}

func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement IsAuthenticated method
	ok := false
	if ok {
		html := `<html>
					<body>
						<form action="/create-playlist" method="get">
							<label for="username">Discogs username:</label>
							<input type="text" id="username" name="username">
							<button type="submit">Create playlist</button>
						</form>
					</body>
				</html>`
		fmt.Fprint(w, html)
	} else {
		html := `<html>
					<body>
						<a href="/login">Log in with Spotify</a>
					</body>
				</html>`
		fmt.Fprint(w, html)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	url := s.oauthController.GetAuthUrl()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter to prevent CSRF
	state := r.FormValue("state")
	code := r.FormValue("code")
	service, err := s.oauthController.GetServiceFromCallback(state, code)
	if err != nil {
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	s.playlistCreator.SetSpotifyService(service)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

		if errors.Cause(err) == spotify.ErrUnauthorized {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
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
