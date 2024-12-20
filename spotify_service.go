package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/martiriera/discogs-spotify/entities"
)

type SpotifyService interface {
	GetAlbumUri(artist string, album string) (string, error)
}

type HttpSpotifyService struct {
	client HttpClient
	token  string
}

func NewHttpSpotifyService(client HttpClient, token string) *HttpSpotifyService {
	return &HttpSpotifyService{client: client, token: token}
}

func (s *HttpSpotifyService) GetAlbumUri(artist string, album string) (string, error) {
	const path = "https://api.spotify.com/v1/search"
	query := url.QueryEscape("album:" + album + " artist:" + artist)
	route := fmt.Sprintf("%s?q=%s&type=album", path, query)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return "", fmt.Errorf("error creating Spotify search request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting Spotify search: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if err := s.setAccessToken(); err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, err = s.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("error requesting Spotify search: %v", err)
		}
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response entities.SpotifySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("unable to parse Spotify response, %v", err)
	}

	if len(response.Albums.Items) == 0 {
		return "not_found", nil
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

	r, _ := http.NewRequest(http.MethodPost, authUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(r)
	if err != nil {
		return fmt.Errorf("error requesting Spotify access token: %v", err)
	}
	defer resp.Body.Close()

	var response entities.SpotifyAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("unable to parse Spotify auth response, %v", err)
	}
	s.token = response.AccessToken
	return nil
}
