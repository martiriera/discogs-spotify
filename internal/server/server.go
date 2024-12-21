package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/util"
	"github.com/pkg/errors"
)

type Server struct {
	playlistCreator *playlist.PlaylistCreator
	http.Handler
}

func NewServer(playlistCreator *playlist.PlaylistCreator) *Server {
	s := new(Server)

	s.playlistCreator = playlistCreator
	router := http.NewServeMux()

	router.Handle("/create-playlist", http.HandlerFunc(s.createPlaylistHandler))

	s.Handler = router
	return s
}

func (s *Server) createPlaylistHandler(w http.ResponseWriter, r *http.Request) {
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
