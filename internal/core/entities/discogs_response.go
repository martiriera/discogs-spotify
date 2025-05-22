package entities

import "strings"

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

func (r *DiscogsListResponse) GetPagination() DiscogsPagination {
	return DiscogsPagination{}
}

func (r *DiscogsListResponse) GetReleases() []DiscogsRelease {
	releases := make([]DiscogsRelease, len(r.Items))
	for i, item := range r.Items {
		titleParts := strings.SplitN(item.DisplayTitle, " - ", 2)
		artist := DiscogsArtist{Name: titleParts[0]}
		title := titleParts[1]
		releases[i] = DiscogsRelease{
			BasicInformation: DiscogsBasicInformation{
				ID:      item.ID,
				Title:   title,
				Artists: []DiscogsArtist{artist},
			},
		}
	}
	return releases
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
	Formats  []DiscogsFormat `json:"formats"`
}

type DiscogsArtist struct {
	Name string `json:"name"`
}

type DiscogsFormat struct {
	Name         string   `json:"name"`
	Quantity     string   `json:"qty"`
	Descriptions []string `json:"descriptions"`
}

type DiscogsListResponse struct {
	CreatedTs   string            `json:"created_ts"`
	ModifiedTs  string            `json:"modified_ts"`
	Name        string            `json:"name"`
	ListID      int               `json:"list_id"`
	URL         string            `json:"url"`
	Items       []DiscogsListItem `json:"items"`
	ResourceURL string            `json:"resource_url"`
	Public      bool              `json:"public"`
	Description string            `json:"description"`
}

type DiscogsListItem struct {
	Comment      string `json:"comment"`
	DisplayTitle string `json:"display_title"`
	URI          string `json:"uri"`
	ImageURL     string `json:"image_url"`
	ResourceURL  string `json:"resource_url"`
	Type         string `json:"type"`
	ID           int    `json:"id"`
}
