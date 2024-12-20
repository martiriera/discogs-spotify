package main

import (
	"log"
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/server"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

func main() {
	creator := playlist.NewPlaylistCreator(
		discogs.NewHttpDiscogsService(&http.Client{}),
		spotify.NewHttpSpotifyService(&http.Client{}, ""),
	)
	s := server.NewServer(creator)

	if err := http.ListenAndServe(":5000", s); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
