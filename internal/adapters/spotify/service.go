package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/martiriera/discogs-spotify/internal/adapters/client"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
	errorWrapper "github.com/martiriera/discogs-spotify/internal/core/errors"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

var (
	ErrSpotifyAPI             = errorWrapper.New("spotify API error")
	ErrSpotifyUnauthorized    = errorWrapper.New("spotify unauthorized error")
	ErrSpotifyInvalidResponse = errorWrapper.New("spotify invalid response error")
)

type HTTPService struct {
	client          client.HTTPClient
	contextProvider ports.ContextPort
}

const basePath = "https://api.spotify.com/v1"

func NewHTTPService(client client.HTTPClient, contextProvider ports.ContextPort) *HTTPService {
	return &HTTPService{
		client:          client,
		contextProvider: contextProvider,
	}
}

func (s *HTTPService) GetAlbumID(ctx context.Context, album entities.Album) (string, error) {
	query := "album:" + album.Title + " artist:" + album.Artist
	if !album.Reissue {
		query += " year:" + strconv.Itoa(album.Year)
	}
	query = url.QueryEscape(query)
	// double encoding
	route := fmt.Sprintf("%s?q=%s&type=album&limit=1", basePath+"/search", url.QueryEscape(query))

	resp, err := doRequest[entities.SpotifySearchResponse](ctx, s, http.MethodGet, route, nil)
	if err != nil {
		return "", err
	}

	if len(resp.Albums.Items) == 0 {
		fmt.Println("no album found for", album.Artist, album.Title)
		return "", nil
	}

	return resp.Albums.Items[0].ID, nil
}

func (s *HTTPService) GetSpotifyUserID(ctx context.Context) (string, error) {
	userID, err := s.contextProvider.GetUserID(ctx)
	if err == nil && userID != "" {
		return userID, nil
	}

	resp, err := doRequest[entities.SpotifyUserResponse](ctx, s, http.MethodGet, basePath+"/me", nil)
	if err != nil {
		return "", err
	}

	if err := s.contextProvider.SetUserID(ctx, resp.ID); err != nil {
		return "", errorWrapper.Wrap(err, "error setting spotify user id in context")
	}

	return resp.ID, nil
}

func (s *HTTPService) CreatePlaylist(ctx context.Context, name string, description string) (entities.SpotifyPlaylist, error) {
	userID, err := s.GetSpotifyUserID(ctx)
	if err != nil {
		return entities.SpotifyPlaylist{}, err
	}

	route := fmt.Sprintf("%s/users/%s/playlists", basePath, userID)

	reqBody := map[string]any{
		"name":        name,
		"description": description,
		"public":      false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return entities.SpotifyPlaylist{}, errorWrapper.Wrap(err, "error marshaling request body")
	}

	resp, err := doRequest[entities.SpotifyPlaylistResponse](ctx, s, http.MethodPost, route, bytes.NewBuffer(jsonBody))
	if err != nil {
		return entities.SpotifyPlaylist{}, err
	}

	return entities.SpotifyPlaylist{
		ID:   resp.ID,
		Name: resp.Name,
		URL:  resp.ExternalURLs.Spotify,
	}, nil
}

func (s *HTTPService) AddToPlaylist(ctx context.Context, playlistID string, uris []string) error {
	route := fmt.Sprintf("%s/playlists/%s/tracks", basePath, playlistID)

	reqBody := map[string]any{
		"uris": uris,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return errorWrapper.Wrap(err, "error marshaling request body")
	}

	_, err = doRequest[map[string]any](ctx, s, http.MethodPost, route, bytes.NewBuffer(jsonBody))
	return err
}

func (s *HTTPService) GetAlbumsTrackUris(ctx context.Context, albums []string) ([]string, error) {
	if len(albums) == 0 {
		return []string{}, nil
	}

	// Spotify API allows a maximum of 20 IDs per request
	batchSize := 20
	var allTrackURIs []string

	for i := 0; i < len(albums); i += batchSize {
		end := i + batchSize
		if end > len(albums) {
			end = len(albums)
		}

		batch := albums[i:end]
		ids := strings.Join(batch, ",")
		route := fmt.Sprintf("%s/albums?ids=%s", basePath, ids)

		resp, err := doRequest[entities.SpotifyAlbumsResponse](ctx, s, http.MethodGet, route, nil)
		if err != nil {
			return nil, err
		}

		for _, album := range resp.Albums {
			for _, track := range album.Tracks.Items {
				allTrackURIs = append(allTrackURIs, track.URI)
			}
		}
	}

	return allTrackURIs, nil
}

func doRequest[T any](ctx context.Context, s *HTTPService, method, route string, body io.Reader) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, method, route, body)
	if err != nil {
		return nil, errorWrapper.Wrap(ErrSpotifyAPI, err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	token, err := s.contextProvider.GetToken(ctx)
	if err != nil {
		return nil, errorWrapper.Wrap(ErrSpotifyUnauthorized, err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errorWrapper.Wrap(ErrSpotifyAPI, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrSpotifyUnauthorized
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("status: %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, errorWrapper.Wrap(ErrSpotifyAPI, errMsg)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errorWrapper.Wrap(ErrSpotifyInvalidResponse, err.Error())
	}

	return &result, nil
}
