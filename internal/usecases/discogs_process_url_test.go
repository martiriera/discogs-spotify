package usecases

import (
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

func TestProcessDiscogsURL(t *testing.T) {
	tests := []struct {
		category  string
		testCases []struct {
			name        string
			url         string
			expected    *entities.ParsedDiscogsURL
			expectError bool
		}
	}{
		{
			category: "Collection URLs",
			testCases: []struct {
				name        string
				url         string
				expected    *entities.ParsedDiscogsURL
				expectError bool
			}{
				{
					name:        "short collection URL",
					url:         "discogs.com/user/digger/collection",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.CollectionType},
					expectError: false,
				},
				{
					name:        "https collection URL",
					url:         "https://www.discogs.com/user/digger/collection",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.CollectionType},
					expectError: false,
				},
				{
					name:        "https collection URL with language code",
					url:         "https://www.discogs.com/es/user/digger/collection",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.CollectionType},
					expectError: false,
				},
				{
					name:        "www collection URL with language code",
					url:         "www.discogs.com/es/user/digger/collection",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.CollectionType},
					expectError: false,
				},
				{
					name:        "https collection URL with different user",
					url:         "https://www.discogs.com/user/johndoe/collection",
					expected:    &entities.ParsedDiscogsURL{ID: "johndoe", Type: entities.CollectionType},
					expectError: false,
				},
				{
					name:        "https collection URL with query parameter",
					url:         "https://www.discogs.com/user/digger/collection?header=1",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.CollectionType},
					expectError: false,
				},
				{
					name:        "collection URL with subdomain",
					url:         "https://m.discogs.com/user/digger/collection",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.CollectionType},
					expectError: false,
				},
			},
		},
		{
			category: "List URLs",
			testCases: []struct {
				name        string
				url         string
				expected    *entities.ParsedDiscogsURL
				expectError bool
			}{
				{
					name:        "https list URL",
					url:         "https://www.discogs.com/lists/MyList/1545836",
					expected:    &entities.ParsedDiscogsURL{ID: "1545836", Type: entities.ListType},
					expectError: false,
				},
				{
					name:        "www list URL",
					url:         "www.discogs.com/lists/MyList/1545836",
					expected:    &entities.ParsedDiscogsURL{ID: "1545836", Type: entities.ListType},
					expectError: false,
				},
				{
					name:        "list URL with language code",
					url:         "www.discogs.com/es/lists/MyList/1545836",
					expected:    &entities.ParsedDiscogsURL{ID: "1545836", Type: entities.ListType},
					expectError: false,
				},
				{
					name:        "list URL with different name",
					url:         "www.discogs.com/lists/FavoriteAlbums/1545836",
					expected:    &entities.ParsedDiscogsURL{ID: "1545836", Type: entities.ListType},
					expectError: false,
				},
			},
		},
		{
			category: "Wantlist URLs",
			testCases: []struct {
				name        string
				url         string
				expected    *entities.ParsedDiscogsURL
				expectError bool
			}{
				{
					name:        "https wantlist URL",
					url:         "https://www.discogs.com/wantlist?user=digger",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.WantlistType},
					expectError: false,
				},
				{
					name:        "www wantlist URL",
					url:         "www.discogs.com/wantlist?user=digger",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.WantlistType},
					expectError: false,
				},
				{
					name:        "short wantlist URL with language code",
					url:         "discogs.com/es/wantlist?user=digger",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.WantlistType},
					expectError: false,
				},
				{
					name:        "wantlist URL with additional parameters",
					url:         "www.discogs.com/wantlist?user=digger&sort=artist",
					expected:    &entities.ParsedDiscogsURL{ID: "digger", Type: entities.WantlistType},
					expectError: false,
				},
			},
		},
		{
			category: "Invalid URLs",
			testCases: []struct {
				name        string
				url         string
				expected    *entities.ParsedDiscogsURL
				expectError bool
			}{
				{
					name:        "empty URL",
					url:         "",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "discogs homepage",
					url:         "https://www.discogs.com",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "incomplete user URL",
					url:         "www.discogs.com/user/digger",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "incomplete collection URL",
					url:         "www.discogs.com/user/digger/collectio",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "wantlist without user parameter",
					url:         "https://www.discogs.com/wantlist?digger",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "invalid lists URL",
					url:         "https://www.discogs.com/user/digger/lists",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "non-discogs URL",
					url:         "test.com",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "non-discogs URL with www",
					url:         "www.test.com",
					expected:    nil,
					expectError: true,
				},
				{
					name:        "malformed URL",
					url:         "http://[::1]:namedport",
					expected:    nil,
					expectError: true,
				},
			},
		},
	}

	for _, category := range tests {
		t.Run(category.category, func(t *testing.T) {
			for _, tc := range category.testCases {
				t.Run(tc.name, func(t *testing.T) {
					got, err := parseDiscogsURL(tc.url)

					if (err != nil) != tc.expectError {
						t.Errorf("parseDiscogsURL(%q) error = %v, expectError = %v",
							tc.url, err, tc.expectError)
					}

					if !reflect.DeepEqual(got, tc.expected) {
						t.Errorf("parseDiscogsURL(%q) = %v, want %v",
							tc.url, got, tc.expected)
					}
				})
			}
		})
	}
}
