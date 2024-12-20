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
			},
		},
	}
}
