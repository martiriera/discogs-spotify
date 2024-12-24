package spotify

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/util"
	"golang.org/x/oauth2"
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
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	tcs := []struct {
		name     string
		request  func(service SpotifyService) (string, error)
		response *http.Response
		want     string
	}{
		{
			name: "should return album uri",
			request: func(service SpotifyService) (string, error) {
				return service.GetAlbumUri(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})
			},
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
			name: "should return empty string as uri when not found",
			request: func(service SpotifyService) (string, error) {
				return service.GetAlbumUri(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})
			},
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
		{
			name: "should return user info",
			request: func(service SpotifyService) (string, error) {
				return service.GetSpotifyUserInfo(ctx)
			},
			response: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
				"display_name": "John Doe",
				"email": "johndoe@example.com",
				"external_urls": {
					"spotify": "string"
				},
				"href": "https://api.spotify.com/v1/users/123",
				"id": "wizzler",
				"images": [
					{
						"url": "https://i.scdn.co/image/ab67616d00001e02ff9ca10b55ce82ae553c8228",
						"height": 300,
						"width": 300
					}
				],
				"type": "user",
				"uri": "spotify:user:wizzler"
				}`)),
			},
			want: "spotify:user:wizzler",
		},
		// {
		// 	name: "should create playlist",
		// 	request: func(service SpotifyService) (string, error) {
		// 		return service.CreatePlaylist([]string{"spotify:album:4JeLdGuCEO9SF9SnFa9LBh"})
		// 	},
		// 	response: &http.Response{
		// 		StatusCode: 201,
		// 		// TODO: more real id
		// 		Body: io.NopCloser(bytes.NewBufferString(`{"id": "123"}`)),
		// 	},
		// 	want: "123",
		// },
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stubClient := &StubSpotifyHttpClient{Responses: []*http.Response{tc.response}}
			service := NewHttpSpotifyService(stubClient)
			response, err := tc.request(service)
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
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	service := NewHttpSpotifyService(stubClient)
	_, err := service.GetAlbumUri(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})

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
			Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Unauthorized"}`)),
		},
	}
	stubClient := &StubSpotifyHttpClient{Responses: stubResponses}
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	service := NewHttpSpotifyService(stubClient)
	_, err := service.GetAlbumUri(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})

	if err == nil {
		t.Errorf("did expect error, got nil")
	}
	if err != ErrUnauthorized {
		t.Errorf("got %v, want %v", err, ErrUnauthorized)
	}
}
