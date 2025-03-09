package entities

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
