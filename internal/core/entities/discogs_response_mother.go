package entities

const (
	discogsArtistDescendents     = "Descendents"
	discogsAlbumDescriptionAlbum = "Album"
	discogsArtistJimCarrollBand  = "The Jim Carroll Band"
)

func MotherTwoDiscogsAlbums() []DiscogsRelease {
	return []DiscogsRelease{
		{
			BasicInformation: DiscogsBasicInformation{
				Title: "Milo Goes to College",
				Artists: []DiscogsArtist{
					{
						Name: discogsArtistDescendents,
					},
				},
				Formats: []DiscogsFormat{
					{
						Descriptions: []string{"LP", discogsAlbumDescriptionAlbum, "Reissue"},
					},
				},
			},
		},
		{
			BasicInformation: DiscogsBasicInformation{
				Title: "Catholic Boy",
				Artists: []DiscogsArtist{
					{
						Name: discogsArtistJimCarrollBand,
					},
				},
				Formats: []DiscogsFormat{
					{
						Descriptions: []string{"LP", discogsAlbumDescriptionAlbum},
					},
				},
			},
		},
	}
}

func MotherNAlbums(n int) []DiscogsRelease {
	albums := []DiscogsRelease{}
	for range n {
		albums = append(albums, DiscogsRelease{
			BasicInformation: DiscogsBasicInformation{
				Title: "Album",
				Artists: []DiscogsArtist{
					{
						Name: "Artist",
					},
				},
			},
		})
	}
	return albums
}
