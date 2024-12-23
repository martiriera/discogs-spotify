package spotify

import (
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
	GetAlbumUri(artist string, title string) (string, error)
	CreatePlaylist(uris []string) (string, error)
	GetSpotifyUserInfo(c *gin.Context) (string, error)
}

type HttpSpotifyService struct {
	client client.HttpClient
}

func NewHttpSpotifyService(client client.HttpClient) *HttpSpotifyService {
	return &HttpSpotifyService{client: client}
}

func (s *HttpSpotifyService) GetAlbumUri(artist string, title string) (string, error) {
	query := url.QueryEscape("album:" + title + " artist:" + artist)
	route := fmt.Sprintf("%s?q=%s&type=album", basePath+"/search", query)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return "", errors.Wrap(ErrSearchRequest, err.Error())
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", errors.Wrap(ErrSearchRequest, err.Error())
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthorized
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", errors.Wrapf(ErrSearchResponse, "status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response entities.SpotifySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", errors.Wrap(ErrSearchResponse, err.Error())
	}

	if len(response.Albums.Items) == 0 {
		return "", nil
	}

	return response.Albums.Items[0].URI, nil
}

func (s *HttpSpotifyService) CreatePlaylist(uris []string) (string, error) {
	return "", nil
}

func (s *HttpSpotifyService) GetSpotifyUserInfo(c *gin.Context) (string, error) {
	token, ok := c.Get(session.SpotifyTokenKey)
	if !ok {
		return "", errors.Wrap(ErrUnauthorized, "no token found")
	}

	req, err := http.NewRequest(http.MethodGet, basePath+"/me", nil)
	if err != nil {
		return "", errors.Wrap(ErrSearchRequest, err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token.(*oauth2.Token).AccessToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify API returned status %d", resp.StatusCode)
	}

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", userInfo), nil
}
