package spotify

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/entities"
)

type SpotifyServiceMock struct {
	Responses []string
	index     int
}

func (m *SpotifyServiceMock) GetAlbumUri(ctx *gin.Context, album entities.Album) (string, error) {
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

func (m *SpotifyServiceMock) GetSpotifyUserInfo(ctx *gin.Context) (string, error) {
	return m.Responses[0], nil
}
