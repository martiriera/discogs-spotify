package session

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// TODO: Decouple from gorilla and gin

var AuthSessionName = "auth-session"
var SpotifyTokenKey = "spotify-token"
var store *sessions.CookieStore

func Init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func GetSession(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

func GetSpotifyToken(r *http.Request) (string, error) {
	authSession, _ := store.Get(r, AuthSessionName)
	token, ok := authSession.Values[SpotifyTokenKey].(string)
	if !ok {
		return "", errors.New("session: token not found")
	}
	return token, nil
}

func SetSpotifyToken(c *gin.Context, token string) error {
	authSession, _ := store.Get(c.Request, AuthSessionName)
	authSession.Values[SpotifyTokenKey] = string(token)
	return authSession.Save(c.Request, c.Writer)
}
