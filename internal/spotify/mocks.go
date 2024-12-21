package spotify

import "net/http"

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

func (m *SpotifyServiceMock) GetSpotifyUserInfo(client *http.Client) (string, error) {
	return m.Responses[0], nil
}
