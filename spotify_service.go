package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SpotifyService interface {
	GetAlbumTitles(url string) ([]string, error)
}

type HttpSpotifyService struct {
	client HTTPClient
}

func (r *HttpSpotifyService) Do(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
}

func NewHttpSpotifyService(client HTTPClient) *HttpSpotifyService {
	return &HttpSpotifyService{client: client}
}

type SpotifyResponse struct {
}

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// func (s *HttpSpotifyService) GetAlbumId(album string, artist string) (string, error) {
// 	const uri = "https://api.spotify.com/v1/search?type=album"
// 	req, err := http.NewRequest(http.MethodGet, uri, nil)
// 	if err != nil {
// 		return "", err
// 	}
// }

func (s *HttpSpotifyService) GetAccessToken() (string, error) {
	const authUrl = "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))

	r, _ := http.NewRequest(http.MethodPost, authUrl, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var response SpotifyAuthResponse
	resp, err := s.client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	return response.AccessToken, nil
}
