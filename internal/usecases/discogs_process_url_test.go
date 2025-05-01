package usecases

import (
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

func TestProcessDiscogsURL(t *testing.T) {
	tcs := []struct {
		name        string
		url         string
		expected    *entities.DiscogsInputURL
		expectError bool
	}{
		{
			"short",
			"discogs.com/user/digger/collection",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.CollectionType},
			false,
		},
		{
			"https",
			"https://www.discogs.com/user/digger/collection",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.CollectionType},
			false,
		},
		{
			"https es",
			"https://www.discogs.com/es/user/digger/collection",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.CollectionType},
			false,
		},
		{
			"www es",
			"www.discogs.com/es/user/digger/collection",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.CollectionType},
			false,
		},
		{
			"https other user",
			"https://www.discogs.com/user/johndoe/collection",
			&entities.DiscogsInputURL{ID: "johndoe", Type: entities.CollectionType},
			false,
		},
		{
			"https with header",
			"https://www.discogs.com/user/digger/collection?header=1",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.CollectionType},
			false,
		},
		{
			"https list",
			"https://www.discogs.com/lists/MyList/1545836",
			&entities.DiscogsInputURL{ID: "1545836", Type: entities.ListType},
			false,
		},
		{
			"www list",
			"www.discogs.com/lists/MyList/1545836",
			&entities.DiscogsInputURL{ID: "1545836", Type: entities.ListType},
			false,
		},
		{
			"https wantlist",
			"https://www.discogs.com/wantlist?user=digger",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.WantlistType},
			false,
		},
		{
			"www wantlist",
			"www.discogs.com/wantlist?user=digger",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.WantlistType},
			false,
		},
		{
			"short wantlist es",
			"discogs.com/es/wantlist?user=digger",
			&entities.DiscogsInputURL{ID: "digger", Type: entities.WantlistType},
			false,
		},
		{
			"wrong",
			"https://www.discogs.com",
			nil,
			true,
		},
		{
			"wrong collection",
			"www.discogs.com/user/digger",
			nil,
			true,
		},
		{
			"incomplete collection",
			"www.discogs.com/user/digger/collectio",
			nil,
			true,
		},
		{
			"wrong query",
			"https://www.discogs.com/wantlist?digger",
			nil,
			true,
		},
		{
			"wrong lists",
			"https://www.discogs.com/user/digger/lists",
			nil,
			true,
		},
		{
			"random url",
			"test.com",
			nil,
			true,
		},
		{
			"www random url",
			"test.com",
			nil,
			true,
		},
	}
	for _, tc := range tcs {
		got, err := parseDiscogsURL(tc.url)
		if (err != nil) != tc.expectError {
			t.Errorf("error = %v, expectError = %v", err, tc.expectError)
		}
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("got %v, want %v: %s", got, tc.expected, tc.name)
		}
	}
}
