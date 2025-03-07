package discogs

import "github.com/martiriera/discogs-spotify/internal/entities"

type DiscogsServiceMock struct {
	Response []entities.DiscogsRelease
	Error    error
}

func (m *DiscogsServiceMock) GetCollectionReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}

func (m *DiscogsServiceMock) GetWantlistReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}

func (m *DiscogsServiceMock) GetListReleases(listID string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}
