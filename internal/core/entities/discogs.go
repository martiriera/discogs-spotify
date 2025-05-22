package entities

type URLType string

func (t URLType) String() string {
	return string(t)
}

const (
	CollectionType URLType = "collection"
	WantlistType   URLType = "wantlist"
	ListType       URLType = "list"
)

type ParsedDiscogsURL struct {
	ID   string
	Type URLType
}
