package discogs

import "github.com/martiriera/discogs-spotify/internal/entities"

type DiscogsServiceMock struct {
	Response []entities.DiscogsRelease
}

func (m *DiscogsServiceMock) GetReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return m.Response, nil
}
