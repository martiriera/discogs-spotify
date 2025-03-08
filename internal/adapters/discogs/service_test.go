package discogs

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

type StubDiscogsHTTPClient struct {
	Responses   []http.Response
	CalledCount int
	Error       error
}

func (s *StubDiscogsHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	if s.CalledCount >= len(s.Responses) {
		return nil, s.Error
	}
	response := (s.Responses)[s.CalledCount]
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

	stubClient := &StubDiscogsHTTPClient{Responses: []http.Response{*stubResponse}}
	service := NewHTTPService(stubClient)

	tcs := []struct {
		name         string
		fn           func(string) ([]entities.DiscogsRelease, error)
		responseBody string
	}{
		{
			name:         "GetCollectionReleases",
			fn:           service.GetCollectionReleases,
			responseBody: generateResponseBody(t, "releases"),
		},
		{
			name:         "GetWantlistReleases",
			fn:           service.GetWantlistReleases,
			responseBody: generateResponseBody(t, "wants"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stubClient.Responses[0].Body = io.NopCloser(bytes.NewBufferString(tc.responseBody))
			stubClient.CalledCount = 0

			response, err := tc.fn("digger")
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
				t.Errorf("got %d, want 1986", response[0].BasicInformation.Year)
			}
			if response[0].BasicInformation.Artists[0].Name != "The Smiths" {
				t.Errorf("got %s, want The Smiths", response[0].BasicInformation.Artists[0].Name)
			}
			if stubClient.CalledCount != 1 {
				t.Errorf("got %d calls, want 1", stubClient.CalledCount)
			}
		})
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
						"last": "https://api.discogs.com/users/digger/collection/folders/0/releases?per_page=1&page=2",
						"next": "https://api.discogs.com/users/digger/collection/folders/0/releases?per_page=1&page=2"
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

	stubClient := &StubDiscogsHTTPClient{Responses: stubResponses}
	service := NewHTTPService(stubClient)
	response, err := service.GetCollectionReleases("digger")
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
	stubClient := &StubDiscogsHTTPClient{Responses: []http.Response{*stubResponse}}
	service := NewHTTPService(stubClient)
	_, err := service.GetCollectionReleases("digger")
	if err == nil {
		t.Errorf("error is nil")
	}
}

func TestDiscogsServiceUnauthorized(t *testing.T) {
	stubResponse := &http.Response{
		StatusCode: 401,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "You must authenticate to access this resource"}`)),
	}
	stubClient := &StubDiscogsHTTPClient{Responses: []http.Response{*stubResponse}}
	service := NewHTTPService(stubClient)
	_, err := service.GetCollectionReleases("digger")
	if err == nil {
		t.Errorf("error is nil")
	}
}

func generateResponseBody(t *testing.T, property string) string {
	t.Helper()
	return `{
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
		"` + property + `": [{
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
	}`
}
