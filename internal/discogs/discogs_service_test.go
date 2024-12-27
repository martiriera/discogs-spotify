package discogs

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type StubDiscogsHttpClient struct {
	Responses   []http.Response
	index       int
	CalledCount int
	Error       error
}

func (s *StubDiscogsHttpClient) Do(req *http.Request) (*http.Response, error) {
	if s.index >= len(s.Responses) {
		return nil, s.Error
	}
	response := (s.Responses)[s.index]
	s.index++
	s.CalledCount++
	return &response, nil
}

func TestDiscogsService(t *testing.T) {
	stubResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`{
				"pagination": {
					"page": 1,
					"pages": 1,
					"per_page": 50,
					"items": 1,
					"urls": {
						"last": "",
						"next": ""
					}
				},
				"releases": [{
					"id": 1,
					"instance_id": 1,
					"date_added": "2021-01-01",
					"basic_information": {
						"id": 1,
						"master_id": 1,
						"title": "The Queen Is Dead",
						"year": 1986,
						"artists": [{
							"name": "The Smiths"
						}]
					}
				}]
			}`)),
	}

	stubClient := &StubDiscogsHttpClient{Responses: []http.Response{*stubResponse}}
	service := NewHttpDiscogsService(stubClient)
	response, err := service.GetReleases("digger")
	if err != nil {
		t.Errorf("did not expect an error, got %v", err)
	}
	if len(response) != 1 {
		t.Errorf("got %d albums, want 1", len(response))
	}
	if response[0].BasicInformation.Title != "The Queen Is Dead" {
		t.Errorf("got %s, want The Queen Is Dead", response[0].BasicInformation.Title)
	}
	if response[0].BasicInformation.Year != 1986 {
		t.Errorf("got %d, want 1986 {", response[0].BasicInformation.Year)
	}
	if response[0].BasicInformation.Artists[0].Name != "The Smiths" {
		t.Errorf("got %s, want The Smiths", response[0].BasicInformation.Artists[0].Name)
	}
	if stubClient.CalledCount != 1 {
		t.Errorf("got %d calls, want 1", stubClient.CalledCount)
	}
}

func TestDiscogsServicePagination(t *testing.T) {
	stubResponses := []http.Response{
		{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(`{
				"pagination": {
					"page": 1,
					"pages": 2,
					"per_page": 1,
					"items": 2,
					"urls": {
						"last": "test.com?page=2",
						"next": "test.com?page=2"
					}
				},
				"releases": [{
					"id": 1,
					"instance_id": 1,
					"date_added": "2021-01-01",
					"basic_information": {
						"id": 1,
						"master_id": 1,
						"title": "The Queen Is Dead",
						"year": 1986,
						"artists": [{
							"name": "The Smiths"
						}]
					}
				}]
			}`)),
		},
		{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(`{
				"pagination": {
					"page": 2,
					"pages": 2,
					"per_page": 1,
					"items": 2,
					"urls": {
						"last": "",
						"next": ""
					}
				},
				"releases": [{
					"id": 2,
					"instance_id": 2,
					"date_added": "2021-01-01",
					"basic_information": {
						"id": 2,
						"master_id": 2,
						"title": "Catholic Boy",
						"year": 1980,
						"artists": [{
							"name": "The Jim Carroll Band"
						}]
					}
				}]
			}`)),
		},
	}

	stubClient := &StubDiscogsHttpClient{Responses: stubResponses}
	service := NewHttpDiscogsService(stubClient)
	response, err := service.GetReleases("digger")
	if err != nil {
		t.Errorf("did not expect an error, got %v", err)
	}
	if len(response) != 2 {
		t.Errorf("got %d albums, want 2", len(response))
	}
	if response[0].BasicInformation.Title != "The Queen Is Dead" {
		t.Errorf("got %s, want The Queen Is Dead", response[0].BasicInformation.Title)
	}
	if response[1].BasicInformation.Title != "Catholic Boy" {
		t.Errorf("got %s, want Catholic Boy", response[1].BasicInformation.Title)
	}
	if stubClient.CalledCount != 2 {
		t.Errorf("got %d calls, want 2", stubClient.CalledCount)
	}
}

func TestDiscogsServiceError(t *testing.T) {
	stubResponse := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Internal server error"}`)),
	}
	stubClient := &StubDiscogsHttpClient{Responses: []http.Response{*stubResponse}}
	service := NewHttpDiscogsService(stubClient)
	_, err := service.GetReleases("digger")
	if err == nil {
		t.Errorf("error is nil")
	}
}

func TestDiscogsServiceUnauthorized(t *testing.T) {
	stubResponse := &http.Response{
		StatusCode: 401,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "You must authenticate to access this resource"}`)),
	}
	stubClient := &StubDiscogsHttpClient{Responses: []http.Response{*stubResponse}}
	service := NewHttpDiscogsService(stubClient)
	_, err := service.GetReleases("digger")
	if err == nil {
		t.Errorf("error is nil")
	}
}
