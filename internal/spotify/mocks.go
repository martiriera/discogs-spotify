package spotify

import "github.com/gin-gonic/gin"

type SpotifyServiceMock struct {
	Responses []string
	index     int
}

func (m *SpotifyServiceMock) GetAlbumUri(artist string, title string) (string, error) {
	if m.index >= len(m.Responses) {
		return "", nil
	}
	response := m.Responses[m.index]
	m.index++
	return response, nil
}

func (m *SpotifyServiceMock) CreatePlaylist(uris []string) (string, error) {
	return m.Responses[0], nil
}

func (m *SpotifyServiceMock) GetSpotifyUserInfo(c *gin.Context) (string, error) {
	return m.Responses[0], nil
}
