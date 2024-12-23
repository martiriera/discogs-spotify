package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var AuthSessionName = "auth-session"
var SpotifyTokenKey = "spotify-token"
var store *sessions.CookieStore

func Init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func GetSession(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}
