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

type DiscogsInputURLType string

func (t DiscogsInputURLType) String() string {
	return string(t)
}

const (
	CollectionType DiscogsInputURLType = "collection"
	WantlistType   DiscogsInputURLType = "wantlist"
	ListType       DiscogsInputURLType = "list"
)

type DiscogsInputURL struct {
	ID   string
	Type DiscogsInputURLType
}
