package main

import (
	"fmt"
	"log"
	"net/http"
)

func triggerUseCase(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	creator := newPlaylistCreator(
		NewHttpDiscogsService(&http.Client{}),
		NewHttpSpotifyService(&http.Client{}, ""),
	)
	uris, err := creator.CreatePlaylist(username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Spotify URIs: %v", uris)
}

func main() {
	http.HandleFunc("/trigger", triggerUseCase)
	fmt.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
