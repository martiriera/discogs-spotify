package main

import (
	"log"
	"net/http"
	"os"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/server"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

func main() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	port := os.Getenv("PORT")

	redirectURL := "http://localhost:" + port + "/callback"

	creator := playlist.NewPlaylistCreator(
		discogs.NewHttpDiscogsService(&http.Client{}),
		spotify.NewHttpSpotifyService(&http.Client{}, ""),
	)

	oauth := spotify.NewOAuthController(
		clientID,
		clientSecret,
		redirectURL,
		[]string{"user-read-private", "user-read-email"}, // Adjust scopes as needed
	)

	s := server.NewServer(creator, oauth)

	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, s); err != nil {
		log.Fatalf("could not listen on port %s: %v", port, err)
	}
}
