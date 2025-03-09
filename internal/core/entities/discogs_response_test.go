package entities

import (
	"reflect"
	"testing"
)

func TestDiscogsListResponse_GetReleases(t *testing.T) {
	tests := []struct {
		name     string
		response DiscogsListResponse
		want     []DiscogsRelease
	}{
		{
			name: "single item",
			response: DiscogsListResponse{
				Items: []DiscogsListItem{
					{
						DisplayTitle: "Artist - Title",
						ID:           1,
					},
				},
			},
			want: []DiscogsRelease{
				{
					BasicInformation: DiscogsBasicInformation{
						ID:    1,
						Title: "Title",
						Artists: []DiscogsArtist{
							{Name: "Artist"},
						},
					},
				},
			},
		},
		{
			name: "multiple items",
			response: DiscogsListResponse{
				Items: []DiscogsListItem{
					{
						DisplayTitle: "Artist1 - Title1",
						ID:           1,
					},
					{
						DisplayTitle: "Artist2 - Title2",
						ID:           2,
					},
				},
			},
			want: []DiscogsRelease{
				{
					BasicInformation: DiscogsBasicInformation{
						ID:    1,
						Title: "Title1",
						Artists: []DiscogsArtist{
							{Name: "Artist1"},
						},
					},
				},
				{
					BasicInformation: DiscogsBasicInformation{
						ID:    2,
						Title: "Title2",
						Artists: []DiscogsArtist{
							{Name: "Artist2"},
						},
					},
				},
			},
		},
		{
			name: "empty items",
			response: DiscogsListResponse{
				Items: []DiscogsListItem{},
			},
			want: []DiscogsRelease{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.GetReleases(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiscogsListResponse.GetReleases() = %v, want %v", got, tt.want)
			}
		})
	}
}
