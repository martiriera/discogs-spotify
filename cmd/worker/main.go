package main

import (
	"context"
	"log"
	"time"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/config"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/container"
)

const (
	workerTimeout     = 30 * time.Minute
	workerPollTimeout = 5 * time.Second
)

func main() {
	log.Println("Starting worker...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		return
	}

	c := container.NewContainer(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
	defer cancel()

	if err := runWorkerTask(ctx, c); err != nil {
		log.Printf("Worker task failed: %v", err)
		return
	}

	log.Println("Worker completed successfully")
}

func runWorkerTask(ctx context.Context, _ *container.Container) error {
	log.Println("Running worker task...")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(workerPollTimeout):
		log.Println("Worker task processing completed")
		return nil
	}
}
