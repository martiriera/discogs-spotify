package session

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type GorillaSession struct {
	store *sessions.CookieStore
}

func NewGorillaSession() *GorillaSession {
	return &GorillaSession{
		store: nil,
	}
}

func (gs *GorillaSession) Init(maxAgeSecs int) {
	gob.Register(&oauth2.Token{})
	gs.store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	gs.store.MaxAge(maxAgeSecs)
}

func (gs *GorillaSession) Get(r *http.Request, sessionName string) (map[any]any, error) {
	session, err := gs.store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}
	return session.Values, nil
}

func (gs *GorillaSession) GetData(r *http.Request, key ContextKey) (any, error) {
	session, err := gs.store.Get(r, AuthSessionName)
	if err != nil {
		return nil, err
	}
	return session.Values[string(key)], nil
}

func (gs *GorillaSession) SetData(r *http.Request, w http.ResponseWriter, key ContextKey, value any) error {
	session, err := gs.store.Get(r, AuthSessionName)
	if err != nil {
		return err
	}
	session.Values[string(key)] = value
	return gs.store.Save(r, w, session)
}
