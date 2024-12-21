package spotify

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type StubSpotifyHttpClient struct {
	Responses []*http.Response
	index     int
	Error     error
}

func (s *StubSpotifyHttpClient) Do(req *http.Request) (*http.Response, error) {
	if s.index >= len(s.Responses) {
		return nil, s.Error
	}
	response := s.Responses[s.index]
	s.index++
	return response, s.Error
}

func TestSpotifyService(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	tcs := []struct {
		test     string
		response *http.Response
		want     string
	}{
		{
			test: "should return album uri",
			response: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
				"albums": {
					"href": "https://api.spotify.com/v1/search?offset=0\u0026limit=20\u0026query=album%3ASpring%20Island%20artist%3ADelta%20Sleep\u0026type=album",
					"items": [
						{
							"album_type": "album",
							"id": "4JeLdGuCEO9SF9SnFa9LBh",
							"name": "Spring Island",
							"uri": "spotify:album:4JeLdGuCEO9SF9SnFa9LBh"
						}
					],
					"limit": 20,
					"next": null,
					"offset": 0,
					"previous": null,
					"total": 1
				}
			}`)),
			},
			want: "spotify:album:4JeLdGuCEO9SF9SnFa9LBh",
		},
		{
			test: "should return empty string as uri when not found",
			response: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
				"albums": {
					"href": "",
					"items": [],
					"limit": 20,
					"next": null,
					"offset": 0,
					"previous": null,
					"total": 1
				}
			}`)),
			},
			want: "",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.test, func(t *testing.T) {
			stubClient := &StubSpotifyHttpClient{Responses: []*http.Response{tc.response}}
			service := NewHttpSpotifyService(stubClient, "test_token")
			response, err := service.GetAlbumUri("Delta Sleep", "Spring Island")
			if err != nil {
				t.Errorf("error is not nil: %v", err)
			}
			if response != tc.want {
				t.Errorf("got %s, want %s", response, tc.want)
			}
		})
	}
}

func TestSpotifyServiceError(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	stubResponse := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Bad Request"}`)),
	}
	stubClient := &StubSpotifyHttpClient{Responses: []*http.Response{stubResponse}}
	service := NewHttpSpotifyService(stubClient, "test_token")
	_, err := service.GetAlbumUri("Delta Sleep", "Spring Island")
	want := `status: 400, body: {"message": "Bad Request"}: spotify search response error`
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != want {
		t.Errorf("got %s, want %s", err.Error(), want)
	}
}

func TestSpotifyServiceUnauthorized(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	stubResponses := []*http.Response{
		{
			StatusCode: 401,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Token has expired"}`)),
		},
		{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(`{
				"access_token": "fresh_token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`)),
		},
		{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewBufferString(`{
				"albums": {
				"items": [
					{
						"uri": "spotify:album:4JeLdGuCEO9SF9SnFa9LBh"
					}
				]
			}
		}`)),
		},
	}
	stubClient := &StubSpotifyHttpClient{Responses: stubResponses}
	service := NewHttpSpotifyService(stubClient, "expired_token")
	response, err := service.GetAlbumUri("Delta Sleep", "Spring Island")
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}
	if response != "spotify:album:4JeLdGuCEO9SF9SnFa9LBh" {
		t.Errorf("got %s, want spotify:album:4JeLdGuCEO9SF9SnFa9LBh", response)
	}
}
