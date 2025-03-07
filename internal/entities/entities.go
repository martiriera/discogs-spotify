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

type DiscogsInputUrlType string

func (t DiscogsInputUrlType) String() string {
	return string(t)
}

const (
	CollectionType DiscogsInputUrlType = "collection"
	WantlistType   DiscogsInputUrlType = "wantlist"
	ListType       DiscogsInputUrlType = "list"
)

type DiscogsInputUrl struct {
	ID   string
	Type DiscogsInputUrlType
}
