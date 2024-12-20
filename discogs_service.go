package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/martiriera/discogs-spotify/entities"
)

type HttpDiscogsService struct {
	client HttpClient
}

const basePath = "https://api.discogs.com"

func (r *HttpDiscogsService) Do(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
}

func NewHttpDiscogsService(client HttpClient) *HttpDiscogsService {
	return &HttpDiscogsService{client: client}
}

func (s *HttpDiscogsService) GetReleases(username string) ([]entities.Release, error) {
	// TODO: Add pagination
	// TODO: Handle private list error
	url := basePath + "/users/" + username + "/collection/folders/0/releases"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating Discogs request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting Discogs releases: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response entities.DiscogsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Discogs response, %v ", err)
	}
	return response.Releases, nil
}
