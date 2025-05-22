package spotify

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
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

type MockContextProvider struct {
	token  *oauth2.Token
	userID string
}

func NewMockContextProvider(token *oauth2.Token, userID string) *MockContextProvider {
	return &MockContextProvider{
		token:  token,
		userID: userID,
	}
}

func (m *MockContextProvider) GetToken(_ context.Context) (*oauth2.Token, error) {
	return m.token, nil
}

func (m *MockContextProvider) GetUserID(_ context.Context) (string, error) {
	return m.userID, nil
}

func (m *MockContextProvider) SetUserID(_ context.Context, userID string) error {
	m.userID = userID
	return nil
}

func TestSearchAlbum(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	tcs := []struct {
		name     string
		request  func(service ports.SpotifyPort) ([]entities.SpotifyAlbumItem, error)
		response *http.Response
		want     []entities.SpotifyAlbumItem
	}{
		{
			name: "should return album id",
			request: func(service ports.SpotifyPort) ([]entities.SpotifyAlbumItem, error) {
				return service.SearchAlbum(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})
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
							"uri": "spotify:album:4JeLdGuCEO9SF9SnFa9LBh",
							"artists": [
								{
									"external_urls": {
										"spotify": "https://open.spotify.com/artist/0kbYTNQb4Pb1rPbbaF0pT5"
									},
									"href": "https://api.spotify.com/v1/artists/0kbYTNQb4Pb1rPbbaF0pT5",
									"id": "0kbYTNQb4Pb1rPbbaF0pT5",
									"name": "Delta Sleep",
									"type": "artist",
									"uri": "spotify:artist:0kbYTNQb4Pb1rPbbaF0pT5"
								}
        			]
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
			want: []entities.SpotifyAlbumItem{
				{
					AlbumType: "album",
					ID:        "4JeLdGuCEO9SF9SnFa9LBh",
					Name:      "Spring Island",
					URI:       "spotify:album:4JeLdGuCEO9SF9SnFa9LBh",
					Artists: []entities.SpotifyAlbumArtist{
						{
							ExternalURLs: entities.SpotifyExternalURLs{
								Spotify: "https://open.spotify.com/artist/0kbYTNQb4Pb1rPbbaF0pT5",
							},
							Href: "https://api.spotify.com/v1/artists/0kbYTNQb4Pb1rPbbaF0pT5",
							ID:   "0kbYTNQb4Pb1rPbbaF0pT5",
							Name: "Delta Sleep",
							Type: "artist",
							URI:  "spotify:artist:0kbYTNQb4Pb1rPbbaF0pT5",
						},
					},
				},
			},
		},
		{
			name: "should return empty string as uri when not found",
			request: func(service ports.SpotifyPort) ([]entities.SpotifyAlbumItem, error) {
				return service.SearchAlbum(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})
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
			want: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stubClient := &StubSpotifyHTTPClient{Responses: []*http.Response{tc.response}}
			contextProvider := NewMockContextProvider(&oauth2.Token{AccessToken: "test"}, "wizzler")
			service := NewHTTPService(stubClient, contextProvider)

			response, err := tc.request(service)

			if err != nil {
				t.Errorf("error is not nil: %v", err)
			}
			if !reflect.DeepEqual(response, tc.want) {
				t.Errorf("got %v, want %v", response, tc.want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	tcs := []struct {
		name     string
		request  func(service ports.SpotifyPort) (string, error)
		response *http.Response
		want     string
	}{
		{
			name: "should return user info",
			request: func(service ports.SpotifyPort) (string, error) {
				return service.GetUserID(ctx)
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
			request: func(service ports.SpotifyPort) (string, error) {
				// set user id in context
				ctx := context.WithValue(ctx, session.SpotifyUserIDKey, "wizzler")
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
			contextProvider := NewMockContextProvider(&oauth2.Token{AccessToken: "test"}, "wizzler")
			service := NewHTTPService(stubClient, contextProvider)
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

func TestCreatePlaylist(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	tcs := []struct {
		name     string
		request  func(service ports.SpotifyPort) (string, error)
		response *http.Response
		want     string
	}{
		{
			name: "should create playlist",
			request: func(service ports.SpotifyPort) (string, error) {
				// set user id in context
				ctx := context.WithValue(ctx, session.SpotifyUserIDKey, "wizzler")
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
			contextProvider := NewMockContextProvider(&oauth2.Token{AccessToken: "test"}, "wizzler")
			service := NewHTTPService(stubClient, contextProvider)
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

func TestGetAlbumsTrackUris(t *testing.T) {
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
	contextProvider := NewMockContextProvider(&oauth2.Token{AccessToken: "test"}, "wizzler")
	service := NewHTTPService(stubClient, contextProvider)
	uris, err := service.GetAlbumsTrackUris(ctx, []string{"spotify:album:1", "spotify:album:2"})
	if err != nil {
		t.Errorf("did not expect error, got %v", err)
	}
	if len(uris) != 5 {
		t.Errorf("got %d uris, want 5", len(uris))
	}
}

func TestServiceError(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test")
	stubResponse := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Bad Request"}`)),
	}
	stubClient := &StubSpotifyHTTPClient{Responses: []*http.Response{stubResponse}}
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	contextProvider := NewMockContextProvider(&oauth2.Token{AccessToken: "test"}, "wizzler")
	service := NewHTTPService(stubClient, contextProvider)
	_, err := service.SearchAlbum(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})

	want := `status: 400, body: {"message": "Bad Request"}: spotify API error`
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != want {
		t.Errorf("got %s, want %s", err.Error(), want)
	}
}

func TestServiceUnauthorized(t *testing.T) {
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

	contextProvider := NewMockContextProvider(&oauth2.Token{AccessToken: "test"}, "wizzler")
	service := NewHTTPService(stubClient, contextProvider)
	_, err := service.SearchAlbum(ctx, entities.Album{Artist: "Delta Sleep", Title: "Spring Island"})

	if err == nil {
		t.Errorf("did expect error, got nil")
	}
	if err != ErrSpotifyUnauthorized {
		t.Errorf("got %v, want %v", err, ErrSpotifyUnauthorized)
	}
}
