package entities

const (
	SpotifyAlbumIDMiloGoesToCollege = "6IU592n49Rn36fpEmg9LIq"
	SpotifyAlbumIDEverythingSucks   = "2jUjrmnCfEEwvC4H2twuTI"
	SpotifyAlbumIDCatholicBoy       = "7a1b2c3d4e5f6g7h8i9j0k"
	SpotifyAlbumIDRooms             = "3z9y8x7w6v5u4t3s2r1q0p"
)

func MotherSpotifyAlbums() []SpotifyAlbumItem {
	return []SpotifyAlbumItem{
		{
			AlbumType: "album",
			ID:        SpotifyAlbumIDMiloGoesToCollege,
			Name:      "Milo Goes to College",
			URI:       "spotify:album:" + SpotifyAlbumIDMiloGoesToCollege,
			Artists: []SpotifyAlbumArtist{
				{
					ExternalURLs: SpotifyExternalURLs{
						Spotify: "https://open.spotify.com/artist/1FGH4Bh7g9W6V4fUcKZWp5",
					},
					Href: "https://api.spotify.com/v1/artists/1FGH4Bh7g9W6V4fUcKZWp5",
					ID:   "1FGH4Bh7g9W6V4fUcKZWp5",
					Name: "Descendents",
					Type: "artist",
					URI:  "spotify:artist:1FGH4Bh7g9W6V4fUcKZWp5",
				},
			},
		},
		{
			AlbumType: "album",
			ID:        SpotifyAlbumIDEverythingSucks,
			Name:      "Everything Sucks",
			URI:       "spotify:album:" + SpotifyAlbumIDEverythingSucks,
			Artists: []SpotifyAlbumArtist{
				{
					ExternalURLs: SpotifyExternalURLs{
						Spotify: "https://open.spotify.com/artist/1FGH4Bh7g9W6V4fUcKZWp5",
					},
					Href: "https://api.spotify.com/v1/artists/1FGH4Bh7g9W6V4fUcKZWp5",
					ID:   "1FGH4Bh7g9W6V4fUcKZWp5",
					Name: "Descendents",
					Type: "artist",
					URI:  "spotify:artist:1FGH4Bh7g9W6V4fUcKZWp5",
				},
			},
		},
		{
			AlbumType: "album",
			ID:        SpotifyAlbumIDCatholicBoy,
			Name:      "Catholic Boy",
			URI:       "spotify:album:" + SpotifyAlbumIDCatholicBoy,
			Artists: []SpotifyAlbumArtist{
				{
					ExternalURLs: SpotifyExternalURLs{
						Spotify: "https://open.spotify.com/artist/7x0hyoQ5ynDBnkvAbJKORj",
					},
					Href: "https://api.spotify.com/v1/artists/7x0hyoQ5ynDBnkvAbJKORj",
					ID:   "7x0hyoQ5ynDBnkvAbJKORj",
					Name: "The Jim Carroll Band",
					Type: "artist",
					URI:  "spotify:artist:7x0hyoQ5ynDBnkvAbJKORj",
				},
			},
		},
		{
			AlbumType: "album",
			ID:        SpotifyAlbumIDRooms,
			Name:      "Rooms",
			URI:       "spotify:album:" + SpotifyAlbumIDRooms,
			Artists: []SpotifyAlbumArtist{
				{
					ExternalURLs: SpotifyExternalURLs{
						Spotify: "https://open.spotify.com/artist/7x0hyoQ5ynDBnkvAbJKORj",
					},
					Href: "https://api.spotify.com/v1/artists/7x0hyoQ5ynDBnkvAbJKORj",
					ID:   "7x0hyoQ5ynDBnkvAbJKORj",
					Name: "The Jim Carroll Band",
					Type: "artist",
					URI:  "spotify:artist:7x0hyoQ5ynDBnkvAbJKORj",
				},
			},
		},
	}
}
