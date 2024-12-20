package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type StubDiscogsHttpClient struct {
	Response *http.Response
	Error    error
}

func (s *StubDiscogsHttpClient) Do(req *http.Request) (*http.Response, error) {
	return s.Response, s.Error
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
						"last": "url_last",
						"next": "url_next"
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

	stubClient := &StubDiscogsHttpClient{Response: stubResponse}
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
}

func TestDiscogsServiceError(t *testing.T) {
	stubResponse := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Internal server error"}`)),
	}
	stubClient := &StubDiscogsHttpClient{Response: stubResponse}
	service := NewHttpDiscogsService(stubClient)
	_, err := service.GetReleases("digger")
	if err == nil {
		t.Errorf("error is nil")
	}
}
