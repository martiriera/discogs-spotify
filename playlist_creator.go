package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/martiriera/discogs-spotify/entities"
)

type PlaylistCreator struct {
	discogsService DiscogsService
	spotifyService SpotifyService
}

func newPlaylistCreator(discogsService DiscogsService, spotifyService SpotifyService) *PlaylistCreator {
	return &PlaylistCreator{
		discogsService: discogsService,
		spotifyService: spotifyService,
	}
}

func (c *PlaylistCreator) CreatePlaylist(discogsUsername string) ([]string, error) {
	releases, err := c.discogsService.GetReleases(discogsUsername)
	if err != nil {
		log.Fatal(err)
	}
	albums := getAlbumsFromReleases(releases)
	spotifyUris := []string{}
	for _, album := range albums {
		uri, err := c.spotifyService.GetAlbumUri(album.Artist, album.Title)
		if err != nil {
			log.Fatal(err)
		}
		spotifyUris = append(spotifyUris, uri)
	}
	fmt.Printf("Spotify URIs: %v", spotifyUris)
	return spotifyUris, nil
}

func getAlbumsFromReleases(releases []entities.DiscogsRelease) []entities.Album {
	albums := []entities.Album{}
	for _, release := range releases {
		album := entities.Album{
			Artist: joinArtists(release.Artists),
			Title:  release.BasicInformation.Title,
		}
		albums = append(albums, album)
	}
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
