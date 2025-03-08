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

type StubSpotifyHTTPClient struct {
	Responses []*http.Response
	index     int
	Error     error
}

func (s *StubSpotifyHTTPClient) Do(_ *http.Request) (*http.Response, error) {
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
		request  func(service Service) (string, error)
		response *http.Response
		want     string
	}{
		{
			name: "should return album id",
			request: func(service Service) (string, error) {
				return service.GetAlbumID(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})
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
			want: "4JeLdGuCEO9SF9SnFa9LBh",
		},
		{
			name: "should return empty string as uri when not found",
			request: func(service Service) (string, error) {
				return service.GetAlbumID(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})
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
			request: func(service Service) (string, error) {
				return service.GetSpotifyUserID(ctx)
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
			want: "wizzler",
		},
		{
			name: "should create playlist",
			request: func(service Service) (string, error) {
				ctx.Set(session.SpotifyUserIDKey, "wizzler")
				playlist, err := service.CreatePlaylist(ctx, "Sunday Playlist", "Rock and Roll")
				return playlist.ID, err
			},
			response: &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewBufferString(`{"id": "6rqhFgbbKwnb9MLmUQDhG6"}`)),
			},
			want: "6rqhFgbbKwnb9MLmUQDhG6",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stubClient := &StubSpotifyHTTPClient{Responses: []*http.Response{tc.response}}
			service := NewHTTPService(stubClient)
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

func TestSpotifyGetAlbumsTrackUris(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})
	stubResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`{
			"albums": [
				{
					"tracks": {
						"items": [
							{
								"uri": "spotify:track:1"
							},
							{
								"uri": "spotify:track:2"
							}
						]
					}
				},
				{
					"tracks": {
						"items": [
							{
								"uri": "spotify:track:3"
							},
							{
								"uri": "spotify:track:4"
							},
							{
								"uri": "spotify:track:5"
							}
						]
					}
				}
			]
		}`)),
	}
	stubClient := &StubSpotifyHTTPClient{Responses: []*http.Response{stubResponse}}
	service := NewHTTPService(stubClient)
	uris, err := service.GetAlbumsTrackUris(ctx, []string{"spotify:album:1", "spotify:album:2"})
	if err != nil {
		t.Errorf("did not expect error, got %v", err)
	}
	if len(uris) != 5 {
		t.Errorf("got %d uris, want 5", len(uris))
	}
}

func TestSpotifyServiceError(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	stubResponse := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Bad Request"}`)),
	}
	stubClient := &StubSpotifyHTTPClient{Responses: []*http.Response{stubResponse}}
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	service := NewHTTPService(stubClient)
	_, err := service.GetAlbumID(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})

	want := `status: 400, body: {"message": "Bad Request"}: spotify API response error`
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
	stubClient := &StubSpotifyHTTPClient{Responses: stubResponses}
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	service := NewHTTPService(stubClient)
	_, err := service.GetAlbumID(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})

	if err == nil {
		t.Errorf("did expect error, got nil")
	}
	if err != ErrUnauthorized {
		t.Errorf("got %v, want %v", err, ErrUnauthorized)
	}
}
