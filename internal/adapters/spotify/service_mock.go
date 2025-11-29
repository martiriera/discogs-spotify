package spotify

import (
	"context"
	"time"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

type ServiceMock struct {
	SearchAlbumResponses [][]entities.SpotifyAlbumItem
	CalledCount          int
	SleepMillis          int
}

func (m *ServiceMock) SearchAlbum(_ context.Context, _ entities.Album) ([]entities.SpotifyAlbumItem, error) {
	if m.CalledCount >= len(m.SearchAlbumResponses) {
		return []entities.SpotifyAlbumItem{}, nil
	}
	response := m.SearchAlbumResponses[m.CalledCount]
	m.CalledCount++
	if m.SleepMillis > 0 {
		time.Sleep(time.Duration(m.SleepMillis) * time.Millisecond)
	}
	return response, nil
}

func (_ *ServiceMock) GetUserID(_ context.Context) (string, error) {
	return "wizzler", nil
}

func (_ *ServiceMock) CreatePlaylist(_ context.Context, _, _ string) (entities.SpotifyPlaylist, error) {
	return entities.SpotifyPlaylist{ID: "6rqhFgbbKwnb9MLmUQDhG6", URL: "https://open.spotify.com/playlist/6rqhFgbbKwnb9MLmUQDhG6"}, nil
}

func (m *ServiceMock) AddToPlaylist(_ context.Context, _ string, _ []string) error {
	m.CalledCount++
	return nil
}

func (_ *ServiceMock) GetAlbumsTrackUris(_ context.Context, _ []string) ([]string, error) {
	return []string{"spotify:track:1", "spotify:track:2"}, nil
}
