package entities

type SpotifySearchResponse struct {
	Albums struct {
		Href     string             `json:"href"`
		Items    []SpotifyAlbumItem `json:"items"`
		Limit    int                `json:"limit"`
		Next     string             `json:"next"`
		Offset   int                `json:"offset"`
		Previous string             `json:"previous"`
		Total    int                `json:"total"`
	} `json:"albums"`
}

type SpotifyAlbumItem struct {
	AlbumType string `json:"album_type"`
	Artists   []struct {
		ExternalUrls SpotifyExternalURLs `json:"external_urls"`
		Href         string              `json:"href"`
		ID           string              `json:"id"`
		Name         string              `json:"name"`
		Type         string              `json:"type"`
		URI          string              `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string            `json:"available_markets"`
	ExternalUrls     SpotifyExternalURLs `json:"external_urls"`
	Href             string              `json:"href"`
	ID               string              `json:"id"`
	Images           []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name                 string `json:"name"`
	ReleaseDate          string `json:"release_date"`
	ReleaseDatePrecision string `json:"release_date_precision"`
	TotalTracks          int    `json:"total_tracks"`
	Type                 string `json:"type"`
	URI                  string `json:"uri"`
}

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifyUserResponse struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalUrls SpotifyExternalURLs `json:"external_urls"`
	Followers    struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

type SpotifyPlaylistResponse struct {
	Collaborative bool                `json:"collaborative"`
	ExternalURLs  SpotifyExternalURLs `json:"external_urls"`
	Href          string              `json:"href"`
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Public        bool                `json:"public"`
	SnapshotID    string              `json:"snapshot_id"`
	Type          string              `json:"type"`
	URI           string              `json:"uri"`
}

type SpotifyAlbumsResponse struct {
	Albums []struct {
		ID     string `json:"id"`
		Tracks struct {
			Items []struct {
				URI string `json:"uri"`
			} `json:"items"`
		} `json:"tracks"`
	} `json:"albums"`
}

type SpotifyExternalURLs struct {
	Spotify string `json:"spotify"`
}

type SpotifySnapshotID struct {
	SnapshotID string `json:"snapshot_id"`
}
