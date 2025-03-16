package ports

import "github.com/martiriera/discogs-spotify/internal/core/entities"

type DiscogsPort interface {
	GetCollectionReleases(username string) ([]entities.DiscogsRelease, error)
	GetWantlistReleases(username string) ([]entities.DiscogsRelease, error)
	GetListReleases(listID string) ([]entities.DiscogsRelease, error)
}
