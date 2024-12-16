package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscogsService interface {
	GetAlbumTitles(url string) ([]string, error)
}

type HttpDiscogsService struct {
	client HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (r *HttpDiscogsService) Do(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
}

func NewHttpDiscogsService(client HTTPClient) *HttpDiscogsService {
	return &HttpDiscogsService{client: client}
}

type Response struct {
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

func (h *HttpDiscogsService) GetAlbumTitles() ([]Release, error) {
	// TODO: Add pagination
	const url = "https://api.discogs.com/users/martireir/collection/folders/0/releases"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to parse discogs releases, %v ", err)
	}
	return response.Releases, nil
}
