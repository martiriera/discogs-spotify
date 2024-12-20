package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	playlistCreator *PlaylistCreator
	http.Handler
}

func NewServer(playlistCreator *PlaylistCreator) *Server {
	s := new(Server)

	s.playlistCreator = playlistCreator
	router := http.NewServeMux()

	router.Handle("/create-playlist", http.HandlerFunc(s.createPlaylistHandler))

	s.Handler = router
	return s
}

func (s *Server) createPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	uris, err := s.playlistCreator.CreatePlaylist(username)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(uris)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(response)
}
