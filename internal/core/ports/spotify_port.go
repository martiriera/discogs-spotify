package ports

import (
	"context"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

type SpotifyPort interface {
	SearchAlbum(ctx context.Context, album entities.Album) (string, error)
	GetSpotifyUserID(ctx context.Context) (string, error)
	CreatePlaylist(ctx context.Context, name string, description string) (entities.SpotifyPlaylist, error)
	AddToPlaylist(ctx context.Context, playlistID string, uris []string) error
	GetAlbumsTrackUris(ctx context.Context, albums []string) ([]string, error)
}
