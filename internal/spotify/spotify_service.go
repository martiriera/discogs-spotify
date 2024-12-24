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
	GetAlbumUri(ctx *gin.Context, album entities.Album) (string, error)
	CreatePlaylist(uris []string) (string, error)
	GetSpotifyUserInfo(ctx *gin.Context) (string, error)
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

	resp, err := s.doRequest(ctx, http.MethodGet, route)
	if err != nil {
		return "", err
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

func (s *HttpSpotifyService) GetSpotifyUserInfo(ctx *gin.Context) (string, error) {
	route := basePath + "/me"

	resp, err := s.doRequest(ctx, http.MethodGet, route)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify API returned status %d", resp.StatusCode)
	}

	var userInfo entities.SpotifyUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", errors.Wrap(ErrSearchResponse, err.Error())
	}

	return fmt.Sprintf("%v", userInfo.URI), nil
}

func (s *HttpSpotifyService) doRequest(ctx *gin.Context, method, route string) (*http.Response, error) {
	token, ok := ctx.Get(session.SpotifyTokenKey)
	if !ok {
		return nil, errors.Wrap(ErrUnauthorized, "no token found")
	}

	req, err := http.NewRequest(method, route, nil)
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

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	return resp, nil
}
