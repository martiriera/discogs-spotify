package playlist

import (
	"fmt"
	"strings"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

type PlaylistCreator struct {
	discogsService discogs.DiscogsService
	spotifyService spotify.SpotifyService
}

func NewPlaylistCreator(discogsService discogs.DiscogsService, spotifyService spotify.SpotifyService) *PlaylistCreator {
	return &PlaylistCreator{
		discogsService: discogsService,
		// TODO: Remove this line
		spotifyService: spotifyService,
	}
}

func (c *PlaylistCreator) SetSpotifyService(service spotify.SpotifyService) {
	c.spotifyService = service
}

func (c *PlaylistCreator) CreatePlaylist(discogsUsername string) ([]string, error) {
	if c.spotifyService == nil {
		return nil, fmt.Errorf("spotify service not set")
	}

	releases, err := c.discogsService.GetReleases(discogsUsername)
	if err != nil {
		return nil, err
	}
	albums := parseAlbumsFromReleases(releases)
	spotifyUris := []string{}
	for _, album := range albums {
		uri, err := c.spotifyService.GetAlbumUri(album.Artist, album.Title)
		if err != nil {
			return nil, err
		}
		spotifyUris = append(spotifyUris, uri)
	}
	fmt.Printf("Spotify URIs: %v", spotifyUris)
	return filterNotFounds(spotifyUris), nil
}

func parseAlbumsFromReleases(releases []entities.DiscogsRelease) []entities.Album {
	albums := []entities.Album{}
	for _, release := range releases {
		album := entities.Album{
			Artist: joinArtists(release.BasicInformation.Artists),
			Title:  strings.TrimSpace(release.BasicInformation.Title),
		}
		albums = append(albums, album)
	}
	// Remove duplicates
	return albums
}

// TODO: Necessary?
func joinArtists(artists []entities.DiscogsArtist) string {
	names := []string{}
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return strings.Join(names, ", ")
}

func filterNotFounds(uris []string) []string {
	filtered := []string{}
	for _, uri := range uris {
		if uri != "" {
			filtered = append(filtered, uri)
		}
	}
	return filtered
}
