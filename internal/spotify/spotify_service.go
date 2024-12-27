package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/client"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var ErrSearchRequest = errors.New("spotify search request error")
var ErrSearchResponse = errors.New("spotify search response error")
var ErrAccessTokenRequest = errors.New("spotify token request error")
var ErrAccessTokenResponse = errors.New("spotify token response error")
var ErrUnauthorized = errors.New("spotify unauthorized error")

const basePath = "https://api.spotify.com/v1"

type SpotifyService interface {
	GetAlbumUri(ctx *gin.Context, album entities.Album) (string, error)
	GetSpotifyUserId(ctx *gin.Context) (string, error)
	CreatePlaylist(ctx *gin.Context, name string, description string) (string, error)
	AddToPlaylist(ctx *gin.Context, playlistId string, uris []string) error
}

type HttpSpotifyService struct {
	client client.HttpClient
}

func NewHttpSpotifyService(client client.HttpClient) *HttpSpotifyService {
	return &HttpSpotifyService{client: client}
}

func (s *HttpSpotifyService) GetAlbumUri(ctx *gin.Context, album entities.Album) (string, error) {
	query := url.QueryEscape("album:" + album.Title + " artist:" + album.Artist)
	route := fmt.Sprintf("%s?q=%s&type=album", basePath+"/search", query)

	resp, err := doRequest[entities.SpotifySearchResponse](s, ctx, http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	if len(resp.Albums.Items) == 0 {
		return "", nil
	}

	return resp.Albums.Items[0].URI, nil
}

func (s *HttpSpotifyService) GetSpotifyUserId(ctx *gin.Context) (string, error) {
	route := basePath + "/me"

	resp, err := doRequest[entities.SpotifyUserResponse](s, ctx, http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", resp.ID), nil
}

func (s *HttpSpotifyService) CreatePlaylist(ctx *gin.Context, name string, description string) (string, error) {
	userId := ctx.GetString(session.SpotifyUserIdKey)
	if userId == "" {
		return "", errors.Wrap(ErrSearchRequest, "no user id found on ctx")
	}

	route := basePath + "/users/" + userId + "/playlists"

	body := map[string]string{
		"name":        name,
		"description": description,
	}

	jsonBody := new(bytes.Buffer)
	if err := json.NewEncoder(jsonBody).Encode(body); err != nil {
		return "", errors.Wrap(ErrSearchRequest, err.Error())
	}

	resp, err := doRequest[entities.SpotifyPlaylistResponse](s, ctx, http.MethodPost, route, jsonBody)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (s *HttpSpotifyService) AddToPlaylist(ctx *gin.Context, playlistId string, uris []string) error {
	route := basePath + "/playlists/" + playlistId + "/tracks"

	body := map[string][]string{
		"uris": uris,
	}

	jsonBody := new(bytes.Buffer)
	if err := json.NewEncoder(jsonBody).Encode(body); err != nil {
		return errors.Wrap(ErrSearchRequest, err.Error())
	}

	_, err := doRequest[entities.SpotifySnapshotId](s, ctx, http.MethodPost, route, jsonBody)
	return err
}

func doRequest[T any](s *HttpSpotifyService, ctx *gin.Context, method, route string, body io.Reader) (*T, error) {
	token, ok := ctx.Get(session.SpotifyTokenKey)
	if !ok {
		return nil, errors.Wrap(ErrUnauthorized, "no token found")
	}

	req, err := http.NewRequest(method, route, body)
	if err != nil {
		return nil, errors.Wrap(ErrSearchRequest, err.Error())
	}

	oauthToken, ok := token.(*oauth2.Token)
	if !ok {
		return nil, errors.Wrap(ErrUnauthorized, "invalid token type")
	}

	req.Header.Set("Authorization", "Bearer "+oauthToken.AccessToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrSearchRequest, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.Wrapf(ErrSearchResponse, "status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(ErrSearchResponse, err.Error())
	}

	return &result, nil
}
