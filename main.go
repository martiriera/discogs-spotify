package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/server"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/util"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("ENV") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("No .env file found")
		}
		gin.SetMode(gin.DebugMode)
	}
	clientID := util.AssertEnvVar("SPOTIFY_CLIENT_ID")
	clientSecret := util.AssertEnvVar("SPOTIFY_CLIENT_SECRET")
	port := util.AssertEnvVar("PORT")
	spotifyAuthRedirectURL := util.AssertEnvVar("SPOTIFY_REDIRECT_URI")

	util.AssertEnvVar("SESSION_KEY")

	session := session.NewGorillaSession()
	session.Init(3600)

	playlistController := playlist.NewPlaylistController(
		discogs.NewHTTPService(&http.Client{}),
		spotify.NewHTTPService(&http.Client{}),
	)

	oauthController := spotify.NewOAuthController(
		clientID,
		clientSecret,
		spotifyAuthRedirectURL,
	)

	userController := spotify.NewUserController(spotify.NewHTTPService(&http.Client{}))

	s := server.NewServer(playlistController, oauthController, userController, session)

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      s,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not listen on port %s: %v", port, err)
	}
}
