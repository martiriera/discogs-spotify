package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/martiriera/discogs-spotify/internal/adapters/client"
	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	"github.com/martiriera/discogs-spotify/internal/adapters/email"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/config"
	"github.com/martiriera/discogs-spotify/internal/usecases"
)

const (
	workerTimeout = 30 * time.Minute
	fromEmail     = "Discogs Digest <hola@martiriera.cat>"
)

func main() {
	log.Println("Starting worker...")

	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		log.Fatalf("Failed to load worker configuration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)

	if err := run(ctx, cfg); err != nil {
		cancel()
		log.Printf("Worker failed: %v", err)
		os.Exit(1)
	}
	cancel()

	log.Println("Worker completed successfully")
}

func run(ctx context.Context, cfg *config.WorkerConfig) error {
	factory := client.NewHTTPClientFactory()

	discogsClient := factory.CreateDiscogsClient(cfg.HTTP.DiscogsTimeout, cfg.HTTP.RetryAttempts, cfg.HTTP.RetryDelay)
	spotifyClient := factory.CreateSpotifyClient(cfg.HTTP.SpotifyTimeout, cfg.HTTP.RetryAttempts, cfg.HTTP.RetryDelay)

	discogsService := discogs.NewHTTPService(discogsClient)

	tokenProvider := spotify.NewClientCredentialsProvider(cfg.SpotifyClientID, cfg.SpotifySecret, spotifyClient)
	contextProvider := spotify.NewWorkerContextProvider(tokenProvider)
	spotifyService := spotify.NewHTTPService(spotifyClient, contextProvider)

	runner := usecases.NewWorkerMonthlyReleases(discogsService, spotifyService)

	// Target: next month
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	targetMonth := nextMonth.Month()
	targetYear := nextMonth.Year()

	log.Printf("Fetching collection from: %s", cfg.DiscogsURL)
	log.Printf("Filtering for month: %s %d", targetMonth, targetYear)

	releases, err := runner.Run(ctx, cfg.DiscogsURL, targetMonth)
	if err != nil {
		return err
	}

	log.Printf("Found %d release(s) for %s", len(releases), targetMonth)

	sender := email.NewResendSender(cfg.ResendAPIKey, fromEmail, cfg.ToEmail)
	if err := sender.SendMonthlyDigest(targetMonth, targetYear, releases); err != nil {
		return err
	}

	log.Printf("Email sent to %s", cfg.ToEmail)
	return nil
}
