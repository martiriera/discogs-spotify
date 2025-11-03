package entities

func MotherTwoDiscogsAlbums() []DiscogsRelease {
	return []DiscogsRelease{
		{
			BasicInformation: DiscogsBasicInformation{
				Title: "Milo Goes to College",
				Artists: []DiscogsArtist{
					{
						Name: "Descendents",
					},
				},
				Formats: []DiscogsFormat{
					{
						Descriptions: []string{"LP", "Album", "Reissue"},
					},
				},
			},
		},
		{
			BasicInformation: DiscogsBasicInformation{
				Title: "Catholic Boy",
				Artists: []DiscogsArtist{
					{
						Name: "The Jim Carroll Band",
					},
				},
				Formats: []DiscogsFormat{
					{
						Descriptions: []string{"LP", "Album"},
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
