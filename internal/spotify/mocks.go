package spotify

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
