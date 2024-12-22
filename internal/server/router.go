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

type ApiRouter struct {
	playlistCreator *playlist.PlaylistCreator
}

func NewApiRouter(c *playlist.PlaylistCreator) *http.ServeMux {
	router := &ApiRouter{playlistCreator: c}
	router.playlistCreator = c

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(router.handleMain))
	mux.Handle("/create-playlist", http.HandlerFunc(router.handlePlaylistCreate))
	return mux
}

func (router *ApiRouter) handleMain(w http.ResponseWriter, r *http.Request) {
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

func (router *ApiRouter) handlePlaylistCreate(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	uris, err := router.playlistCreator.CreatePlaylist(username)
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
