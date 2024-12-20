package entities

type DiscogsResponse struct {
	Pagination struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
		Items   int `json:"items"`
		Urls    struct {
			Last string `json:"last"`
			Next string `json:"next"`
		} `json:"urls"`
	} `json:"pagination"`
	Releases []DiscogsRelease `json:"releases"`
}

type DiscogsRelease struct {
	ID               int                     `json:"id"`
	InstanceID       int                     `json:"instance_id"`
	DateAdded        string                  `json:"date_added"`
	BasicInformation DiscogsBasicInformation `json:"basic_information"`
	Artists          []DiscogsArtist         `json:"artists"`
}

type DiscogsBasicInformation struct {
	ID       int    `json:"id"`
	MasterID int    `json:"master_id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
}

type DiscogsArtist struct {
	Name string `json:"name"`
}
