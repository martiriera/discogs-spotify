package discogs

import "github.com/martiriera/discogs-spotify/internal/entities"

type ServiceMock struct {
	Response []entities.DiscogsRelease
	Error    error
}

func (m *ServiceMock) GetCollectionReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}

func (m *ServiceMock) GetWantlistReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}

func (m *ServiceMock) GetListReleases(listID string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}
