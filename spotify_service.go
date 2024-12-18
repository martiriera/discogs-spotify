package main

import (
	"encoding/json"
	"fmt"
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

func (s *HttpSpotifyService) GetAlbumId(album string, artist string) (*SpotifyResponse, error) {
	const path = "https://api.spotify.com/v1/search"
	query := "album:" + album + " artist:" + artist + " type:album"
	route := path + "?q=" + url.QueryEscape(query)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}
	token, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("error getting Spotify access token: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting Spotify album id: %v", err)
	}
	defer resp.Body.Close()
	var response SpotifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("unable to parse Spotify response, %v ", err)
	}
	return &response, nil
}

func (s *HttpSpotifyService) getAccessToken() (string, error) {
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
		// TODO: explain error
		return "", err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// TODO: explain error
		return "", err
	}
	return response.AccessToken, nil
}
