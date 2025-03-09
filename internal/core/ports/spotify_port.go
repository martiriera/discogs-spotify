package ports

import (
	"github.com/gin-gonic/gin"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

type SpotifyPort interface {
	GetAlbumID(ctx *gin.Context, album entities.Album) (string, error)
	GetSpotifyUserID(ctx *gin.Context) (string, error)
	CreatePlaylist(ctx *gin.Context, name string, description string) (entities.SpotifyPlaylist, error)
	AddToPlaylist(ctx *gin.Context, playlistID string, uris []string) error
	GetAlbumsTrackUris(ctx *gin.Context, albums []string) ([]string, error)
}
