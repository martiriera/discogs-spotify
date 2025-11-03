package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

	spotifyClientID := getRequiredEnv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := getRequiredEnv("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURI := getRequiredEnv("SPOTIFY_REDIRECT_URI")
	spotifyProxyURL := getEnvWithDefault("SPOTIFY_PROXY_URL", "")
	sessionKey := getRequiredEnv("SESSION_KEY")

	port := getEnvWithDefault("PORT", "8080")
	env := getEnvWithDefault("ENV", "development")
	sessionMaxAge := getEnvAsIntWithDefault("SESSION_MAX_AGE", 3600)

	discogsTimeout := getEnvAsDurationWithDefault("DISCOGS_TIMEOUT", 30*time.Second)
	spotifyTimeout := getEnvAsDurationWithDefault("SPOTIFY_TIMEOUT", 60*time.Second)
	retryAttempts := getEnvAsIntWithDefault("HTTP_RETRY_ATTEMPTS", 3)
	retryDelay := getEnvAsDurationWithDefault("HTTP_RETRY_DELAY", 1*time.Second)

	readTimeout := getEnvAsDurationWithDefault("SERVER_READ_TIMEOUT", 10*time.Second)
	writeTimeout := getEnvAsDurationWithDefault("SERVER_WRITE_TIMEOUT", 120*time.Second)
	idleTimeout := getEnvAsDurationWithDefault("SERVER_IDLE_TIMEOUT", 120*time.Second)

	return &Config{
		Environment: env,
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
			UseProxy:     env == "development" && spotifyProxyURL != "",
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

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
