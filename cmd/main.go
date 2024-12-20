package main

import (
	"log"
	"net/http"
)

func main() {
	creator := newPlaylistCreator(
		NewHttpDiscogsService(&http.Client{}),
		NewHttpSpotifyService(&http.Client{}, ""),
	)
	server := NewServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
