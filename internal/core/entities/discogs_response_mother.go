package entities

func MotherTwoAlbums() []DiscogsRelease {
	return []DiscogsRelease{
		{
			BasicInformation: DiscogsBasicInformation{
				Title: "Tim",
				Artists: []DiscogsArtist{
					{
						Name: "The Replacements",
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
	for i := 0; i < n; i++ {
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
