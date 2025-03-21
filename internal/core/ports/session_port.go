package ports

import (
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

type SessionPort interface {
	Init(maxAgeSecs int)
	Get(r *http.Request, sessionName string) (map[any]any, error)
	GetData(r *http.Request, key session.ContextKey) (any, error)
	SetData(r *http.Request, w http.ResponseWriter, key session.ContextKey, value any) error
}
