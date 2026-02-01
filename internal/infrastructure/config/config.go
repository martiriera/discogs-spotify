package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/martiriera/discogs-spotify/internal/utils/env"
)

const (
	defaultSessionMaxAge      = 3600 // 1 hour
	defaultDiscogsTimeout     = 30   // 30 seconds
	defaultSpotifyTimeout     = 60   // 60 seconds
	defaultServerReadTimeout  = 10   // 10 seconds
	defaultServerWriteTimeout = 120  // 120 seconds (2 minutes)
	defaultServerIdleTimeout  = 120  // 120 seconds (2 minutes)
)

type Config struct {
	Environment string
	Server      ServerConfig
	Spotify     SpotifyConfig
	Session     SessionConfig
	HTTP        HTTPConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	ProxyURL     string // Auth proxy URL for development
	UseProxy     bool
}

type SessionConfig struct {
	Key       string
	MaxAgeSec int
}

type HTTPConfig struct {
	DiscogsTimeout time.Duration
	SpotifyTimeout time.Duration
	RetryAttempts  int
	RetryDelay     time.Duration
}

func LoadConfig() (*Config, error) {
	// Load .env file if ENV is not set
	if os.Getenv("ENV") == "" {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("no .env file found: %w", err)
		}
	}

	spotifyClientID := env.GetRequired("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := env.GetRequired("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURI := env.GetRequired("SPOTIFY_REDIRECT_URI")
	spotifyProxyURL := env.GetWithDefault("SPOTIFY_PROXY_URL", "")
	sessionKey := env.GetRequired("SESSION_KEY")

	port := env.GetWithDefault("PORT", "8080")
	environment := env.GetWithDefault("ENV", "development")
	sessionMaxAge := env.GetAsIntWithDefault("SESSION_MAX_AGE", defaultSessionMaxAge)

	discogsTimeout := env.GetAsDurationWithDefault("DISCOGS_TIMEOUT", defaultDiscogsTimeout*time.Second)
	spotifyTimeout := env.GetAsDurationWithDefault("SPOTIFY_TIMEOUT", defaultSpotifyTimeout*time.Second)
	retryAttempts := env.GetAsIntWithDefault("HTTP_RETRY_ATTEMPTS", 3)
	retryDelay := env.GetAsDurationWithDefault("HTTP_RETRY_DELAY", 1*time.Second)

	readTimeout := env.GetAsDurationWithDefault("SERVER_READ_TIMEOUT", defaultServerReadTimeout*time.Second)
	writeTimeout := env.GetAsDurationWithDefault("SERVER_WRITE_TIMEOUT", defaultServerWriteTimeout*time.Second)
	idleTimeout := env.GetAsDurationWithDefault("SERVER_IDLE_TIMEOUT", defaultServerIdleTimeout*time.Second)

	return &Config{
		Environment: environment,
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		Spotify: SpotifyConfig{
			ClientID:     spotifyClientID,
			ClientSecret: spotifyClientSecret,
			RedirectURI:  spotifyRedirectURI,
			ProxyURL:     spotifyProxyURL,
			UseProxy:     environment == "development" && spotifyProxyURL != "",
		},
		Session: SessionConfig{
			Key:       sessionKey,
			MaxAgeSec: sessionMaxAge,
		},
		HTTP: HTTPConfig{
			DiscogsTimeout: discogsTimeout,
			SpotifyTimeout: spotifyTimeout,
			RetryAttempts:  retryAttempts,
			RetryDelay:     retryDelay,
		},
	}, nil
}
