package entities

type DiscogsResponse interface {
	GetPagination() DiscogsPagination
	GetReleases() []DiscogsRelease
}

func (r *DiscogsCollectionResponse) GetPagination() DiscogsPagination {
	return r.Pagination
}

func (r *DiscogsCollectionResponse) GetReleases() []DiscogsRelease {
	return r.Releases
}

func (r *DiscogsWantlistResponse) GetPagination() DiscogsPagination {
	return r.Pagination
}

func (r *DiscogsWantlistResponse) GetReleases() []DiscogsRelease {
	return r.Wants
}

type DiscogsCollectionResponse struct {
	Pagination DiscogsPagination `json:"pagination"`
	Releases   []DiscogsRelease  `json:"releases"`
}

type DiscogsWantlistResponse struct {
	Pagination DiscogsPagination `json:"pagination"`
	Wants      []DiscogsRelease  `json:"wants"`
}

type DiscogsPagination struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"per_page"`
	Items   int `json:"items"`
	Urls    struct {
		Last string `json:"last"`
		Next string `json:"next"`
	} `json:"urls"`
}

type DiscogsRelease struct {
	ID               int                     `json:"id"`
	InstanceID       int                     `json:"instance_id"`
	DateAdded        string                  `json:"date_added"`
	BasicInformation DiscogsBasicInformation `json:"basic_information"`
}

type DiscogsBasicInformation struct {
	ID       int             `json:"id"`
	MasterID int             `json:"master_id"`
	Title    string          `json:"title"`
	Year     int             `json:"year"`
	Artists  []DiscogsArtist `json:"artists"`
}

type DiscogsArtist struct {
	Name string `json:"name"`
}
