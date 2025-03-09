package discogs

import "github.com/martiriera/discogs-spotify/internal/core/entities"

type ServiceMock struct {
	Response []entities.DiscogsRelease
	Error    error
}

func (m *ServiceMock) GetCollectionReleases(_ string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}

func (m *ServiceMock) GetWantlistReleases(_ string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}

func (m *ServiceMock) GetListReleases(_ string) ([]entities.DiscogsRelease, error) {
	return m.Response, m.Error
}
