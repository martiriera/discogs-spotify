package session

import (
	"net/http"
)

var AuthSessionName = "auth-session"
var SpotifyTokenKey = "spotify-token"

type Session interface {
	Init()
	Get(r *http.Request, sessionName string) (*SessionData, error)
	GetData(r *http.Request, key string) (interface{}, error)
	SetData(r *http.Request, w http.ResponseWriter, key string, value interface{}) error
}

type SessionData struct {
	Values map[interface{}]interface{}
}
