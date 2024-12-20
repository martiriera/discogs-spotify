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
	Releases []Release `json:"releases"`
}

type Release struct {
	ID               int    `json:"id"`
	InstanceID       int    `json:"instance_id"`
	DateAdded        string `json:"date_added"`
	BasicInformation struct {
		ID       int    `json:"id"`
		MasterID int    `json:"master_id"`
		Title    string `json:"title"`
		Year     int    `json:"year"`
	} `json:"basic_information"`
	Artists []struct {
		Name string `json:"name"`
	}
}
