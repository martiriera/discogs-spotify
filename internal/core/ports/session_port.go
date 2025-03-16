package ports

import (
	"net/http"
)

type SessionPort interface {
	Init(maxAgeSecs int)
	Get(r *http.Request, sessionName string) (map[any]any, error)
	GetData(r *http.Request, key string) (any, error)
	SetData(r *http.Request, w http.ResponseWriter, key string, value any) error
}
