package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type StubHTTPClient struct {
	Response *http.Response
	Error    error
}

func (s *StubHTTPClient) Do(req *http.Request) (*http.Response, error) {
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
						"title": "Album Title",
						"year": 2021
					},
					"artists": [{
						"name": "Artist Name"
					}]
				}]
			}`)),
	}

	stubClient := &StubHTTPClient{Response: stubResponse}
	service := NewHttpDiscogsService(stubClient)
	response, err := service.GetReleases()
	if err != nil {
		t.Errorf("error is not nil")
	}
	if len(response) != 1 {
		t.Errorf("got %d albums, want 1", len(response))
	}
	if response[0].BasicInformation.Title != "Album Title" {
		t.Errorf("got %s, want Album Title", response[0].BasicInformation.Title)
	}
	if response[0].BasicInformation.Year != 2021 {
		t.Errorf("got %d, want 2021", response[0].BasicInformation.Year)
	}
	if response[0].Artists[0].Name != "Artist Name" {
		t.Errorf("got %s, want Artist Name", response[0].Artists[0].Name)
	}
}

func TestDiscogsServiceError(t *testing.T) {
	stubClient := &StubHTTPClient{Error: fmt.Errorf("request error")}
	service := NewHttpDiscogsService(stubClient)
	_, err := service.GetReleases()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDiscogsServiceInvalidJSON(t *testing.T) {
	stubResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
	}

	stubClient := &StubHTTPClient{Response: stubResponse}
	service := NewHttpDiscogsService(stubClient)
	_, err := service.GetReleases()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
