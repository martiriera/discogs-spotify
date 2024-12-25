package session

import (
	"net/http"
)

var AuthSessionName = "auth-session"
var SpotifyTokenKey = "spotify-token"
var SpotifyUserIdKey = "spotify-user-id"

type Session interface {
	Init()
	Get(r *http.Request, sessionName string) (map[any]any, error)
	GetData(r *http.Request, key string) (any, error)
	SetData(r *http.Request, w http.ResponseWriter, key string, value any) error
}
