package main

import (
	"log"
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/server"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("No .env file found")
	}
	clientID := util.AssertEnvVar("SPOTIFY_CLIENT_ID")
	clientSecret := util.AssertEnvVar("SPOTIFY_CLIENT_SECRET")
	port := util.AssertEnvVar("PORT")
	util.AssertEnvVar("SESSION_KEY")

	oauthRedirectUrl := "http://localhost:" + port + "/auth/callback"

	session := session.NewGorillaSession()
	session.Init()

	creator := playlist.NewPlaylistController(
		discogs.NewHttpDiscogsService(&http.Client{}),
		spotify.NewHttpSpotifyService(&http.Client{}),
	)

	oauth := spotify.NewOAuthController(
		clientID,
		clientSecret,
		oauthRedirectUrl,
	)

	s := server.NewServer(creator, oauth, session)

	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, s); err != nil {
		log.Fatalf("could not listen on port %s: %v", port, err)
	}
}
