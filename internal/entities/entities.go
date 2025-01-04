package entities

type Album struct {
	Artist string
	Title  string
}

type SpotifyPlaylist struct {
	ID  string
	URL string
}

type Playlist struct {
	SpotifyPlaylist
	DiscogsReleases int
	SpotifyAlbums   int
}
