package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/martiriera/discogs-spotify/entities"
)

type HttpDiscogsService struct {
	client HttpClient
}

func (r *HttpDiscogsService) Do(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
}

func NewHttpDiscogsService(client HttpClient) *HttpDiscogsService {
	return &HttpDiscogsService{client: client}
}

func (s *HttpDiscogsService) GetReleases() ([]entities.Release, error) {
	// TODO: Add pagination
	const url = "https://api.discogs.com/users/martireir/collection/folders/0/releases"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating Discogs request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting Discogs releases: %v", err)
	}
	defer resp.Body.Close()

	var response entities.DiscogsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Discogs response, %v ", err)
	}
	return response.Releases, nil
}
