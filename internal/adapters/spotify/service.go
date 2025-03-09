package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/adapters/client"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

var ErrRequest = errors.New("spotify API request error")
var ErrResponse = errors.New("spotify API response error")
var ErrUnauthorized = errors.New("spotify API unauthorized error")

const basePath = "https://api.spotify.com/v1"

type HTTPService struct {
	client client.HTTPClient
}

func NewHTTPService(client client.HTTPClient) *HTTPService {
	return &HTTPService{client: client}
}

func (s *HTTPService) GetAlbumID(ctx *gin.Context, album entities.Album) (string, error) {
	query := url.QueryEscape("album:" + album.Title + " artist:" + album.Artist)
	route := fmt.Sprintf("%s?q=%s&type=album&limit=1", basePath+"/search", query)

	resp, err := doRequest[entities.SpotifySearchResponse](s, ctx, http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	if len(resp.Albums.Items) == 0 {
		fmt.Println("no album found for", album.Artist, album.Title)
		return "", nil
	}

	return resp.Albums.Items[0].ID, nil
}

func (s *HTTPService) GetSpotifyUserID(ctx *gin.Context) (string, error) {
	route := basePath + "/me"

	resp, err := doRequest[entities.SpotifyUserResponse](s, ctx, http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", resp.ID), nil
}

func (s *HTTPService) CreatePlaylist(ctx *gin.Context, name string, description string) (entities.SpotifyPlaylist, error) {
	userID := ctx.GetString(session.SpotifyUserIDKey)
	if userID == "" {
		return entities.SpotifyPlaylist{}, errors.Wrap(ErrRequest, "no user id found on ctx")
	}

	route := basePath + "/users/" + userID + "/playlists"

	body := map[string]string{
		"name":        name,
		"description": description,
	}

	jsonBody := new(bytes.Buffer)
	if err := json.NewEncoder(jsonBody).Encode(body); err != nil {
		return entities.SpotifyPlaylist{}, errors.Wrap(ErrRequest, err.Error())
	}

	resp, err := doRequest[entities.SpotifyPlaylistResponse](s, ctx, http.MethodPost, route, jsonBody)
	if err != nil {
		return entities.SpotifyPlaylist{}, err
	}

	return entities.SpotifyPlaylist{ID: resp.ID, URL: resp.ExternalUrls.Spotify}, nil
}

func (s *HTTPService) AddToPlaylist(ctx *gin.Context, playlistID string, uris []string) error {
	route := basePath + "/playlists/" + playlistID + "/tracks"

	body := map[string][]string{
		"uris": uris,
	}

	jsonBody := new(bytes.Buffer)
	if err := json.NewEncoder(jsonBody).Encode(body); err != nil {
		return errors.Wrap(ErrRequest, err.Error())
	}

	_, err := doRequest[entities.SpotifySnapshotID](s, ctx, http.MethodPost, route, jsonBody)
	return err
}

func (s *HTTPService) GetAlbumsTrackUris(ctx *gin.Context, albums []string) ([]string, error) {
	query := url.Values{}
	query.Set("ids", strings.Join(albums, ","))
	route := basePath + "/albums" + "?" + query.Encode()

	resp, err := doRequest[entities.SpotifyAlbumsResponse](s, ctx, http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	var tracks []string
	for _, album := range resp.Albums {
		for _, track := range album.Tracks.Items {
			tracks = append(tracks, track.URI)
		}
	}

	return tracks, nil
}

func doRequest[T any](s *HTTPService, ctx *gin.Context, method, route string, body io.Reader) (*T, error) {
	token, ok := ctx.Get(session.SpotifyTokenKey)
	if !ok {
		return nil, errors.Wrap(ErrUnauthorized, "no token found")
	}

	req, err := http.NewRequest(method, route, body)
	if err != nil {
		return nil, errors.Wrap(ErrRequest, err.Error())
	}

	oauthToken, ok := token.(*oauth2.Token)
	if !ok {
		return nil, errors.Wrap(ErrUnauthorized, "invalid token type")
	}

	req.Header.Set("Authorization", "Bearer "+oauthToken.AccessToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrRequest, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.Wrapf(ErrResponse, "status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(ErrResponse, err.Error())
	}

	return &result, nil
}
