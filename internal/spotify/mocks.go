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

func (m *SpotifyServiceMock) GetSpotifyUserId(ctx *gin.Context) (string, error) {
	return "wizzler", nil
}

func (m *SpotifyServiceMock) CreatePlaylist(ctx *gin.Context, name string, description string) (string, error) {
	return "6rqhFgbbKwnb9MLmUQDhG6", nil
}

func (m *SpotifyServiceMock) AddToPlaylist(ctx *gin.Context, playlistId string, uris []string) error {
	return nil
}
