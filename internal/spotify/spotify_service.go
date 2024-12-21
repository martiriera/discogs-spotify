package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/martiriera/discogs-spotify/internal/client"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/pkg/errors"
)

var ErrSearchRequest = errors.New("spotify search request error")
var ErrSearchResponse = errors.New("spotify search response error")
var ErrAccessTokenRequest = errors.New("spotify token request error")
var ErrAccessTokenResponse = errors.New("spotify token response error")

const basePath = "https://api.spotify.com/v1"

type SpotifyService interface {
	GetAlbumUri(artist string, title string) (string, error)
	CreatePlaylist(uris []string) (string, error)
	GetSpotifyUserInfo(client *http.Client) (string, error)
}

type HttpSpotifyService struct {
	client client.HttpClient
	token  string
}

func NewHttpSpotifyService(client client.HttpClient, token string) *HttpSpotifyService {
	return &HttpSpotifyService{client: client, token: token}
}

func (s *HttpSpotifyService) GetAlbumUri(artist string, title string) (string, error) {
	query := url.QueryEscape("album:" + title + " artist:" + artist)
	route := fmt.Sprintf("%s?q=%s&type=album", basePath+"/search", query)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return "", errors.Wrap(ErrSearchRequest, err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", errors.Wrap(ErrSearchRequest, err.Error())
	}

	if s.token == "" || resp.StatusCode == http.StatusUnauthorized {
		if err := s.setAccessToken(); err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, err = s.client.Do(req)
		if err != nil {
			return "", errors.Wrap(ErrSearchRequest, err.Error())
		}
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

func (s *HttpSpotifyService) GetSpotifyUserInfo(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Spotify API returned status %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", userInfo), nil
}

func (s *HttpSpotifyService) setAccessToken() error {
	const authUrl = "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_ID or SPOTIFY_CLIENT_SECRET environment variable not set")
	}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	r, err := http.NewRequest(http.MethodPost, authUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(ErrAccessTokenRequest, err.Error())
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(r)
	if err != nil {
		return errors.Wrap(ErrAccessTokenRequest, err.Error())
	}
	defer resp.Body.Close()

	var response entities.SpotifyAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(ErrAccessTokenResponse, err.Error())
	}
	s.token = response.AccessToken
	return nil
}
