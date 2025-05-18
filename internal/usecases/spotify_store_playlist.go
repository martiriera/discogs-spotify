package usecases

import (
	"context"
	"fmt"
	"sync"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/database"
)

type StorePlaylistUseCase struct {
	spotifyService ports.SpotifyPort
	repository     *database.PlaylistRepository
}

func NewStorePlaylistUseCase(spotifyService ports.SpotifyPort, repository *database.PlaylistRepository) *StorePlaylistUseCase {
	return &StorePlaylistUseCase{
		spotifyService: spotifyService,
		repository:     repository,
	}
}

func (u *StorePlaylistUseCase) Execute(ctx context.Context, albums []entities.Album) error {
	// Convert to our database model
	var dbAlbums []database.PlaylistAlbum
	for _, item := range albums {
		track := database.PlaylistAlbum{
			AlbumName:   item.Title,
			ArtistName:  item.Artist,
			ReleaseDate: item.ReleaseDate,
		}
		dbAlbums = append(dbAlbums, track)
	}

	// Store in database
	return u.repository.StorePlaylistAlbums(dbAlbums)
}

// ExecuteAsync runs the store operation in a goroutine
func (uc *StorePlaylistUseCase) ExecuteAsync(ctx context.Context, albums []entities.Album, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := uc.Execute(ctx, albums)
		if err != nil {
			fmt.Println("Error storing playlist albums:", err)
			return
		}
		fmt.Println("Playlist albums stored successfully")
	}()
}
