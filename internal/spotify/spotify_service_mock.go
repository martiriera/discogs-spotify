package spotify

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/entities"
)

type SpotifyServiceMock struct {
	Responses   []string
	CalledCount int
}

func (m *SpotifyServiceMock) GetAlbumId(ctx *gin.Context, album entities.Album) (string, error) {
	if m.CalledCount >= len(m.Responses) {
		return "", nil
	}
	response := m.Responses[m.CalledCount]
	m.CalledCount++
	return response, nil
}

func (m *SpotifyServiceMock) GetSpotifyUserId(ctx *gin.Context) (string, error) {
	return "wizzler", nil
}

func (m *SpotifyServiceMock) CreatePlaylist(ctx *gin.Context, name string, description string) (entities.Playlist, error) {
	return entities.Playlist{ID: "6rqhFgbbKwnb9MLmUQDhG6", URL: "https://open.spotify.com/playlist/6rqhFgbbKwnb9MLmUQDhG6"}, nil
}

func (m *SpotifyServiceMock) AddToPlaylist(ctx *gin.Context, playlistId string, uris []string) error {
	m.CalledCount++
	return nil
}

func (m *SpotifyServiceMock) GetAlbumsTrackUris(ctx *gin.Context, albums []string) ([]string, error) {
	return []string{"spotify:track:1", "spotify:track:2"}, nil
}
