package ports

import (
	"context"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

type DiscogsPort interface {
	GetCollectionReleases(ctx context.Context, username string) ([]entities.DiscogsRelease, error)
	GetWantlistReleases(ctx context.Context, username string) ([]entities.DiscogsRelease, error)
	GetListReleases(ctx context.Context, listID string) ([]entities.DiscogsRelease, error)
}
