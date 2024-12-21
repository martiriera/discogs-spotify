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

type SpotifyService interface {
	GetAlbumUri(artist string, title string) (string, error)
}

type HttpSpotifyService struct {
	client client.HttpClient
	token  string
}

func NewHttpSpotifyService(client client.HttpClient, token string) *HttpSpotifyService {
	return &HttpSpotifyService{client: client, token: token}
}

func (s *HttpSpotifyService) GetAlbumUri(artist string, title string) (string, error) {
	const path = "https://api.spotify.com/v1/search"
	query := url.QueryEscape("album:" + title + " artist:" + artist)
	route := fmt.Sprintf("%s?q=%s&type=album", path, query)
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
