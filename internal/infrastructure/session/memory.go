package session

import (
	"errors"
	"net/http"
	"time"
)

type InMemorySession struct {
	Data      map[any]any
	ExpiresAt int
}

func NewInMemorySession() *InMemorySession {
	return &InMemorySession{
		Data: make(map[any]any),
	}
}

func (s *InMemorySession) Init(maxAgeSecs int) {
	s.ExpiresAt = int(time.Now().Add(time.Duration(maxAgeSecs) * time.Second).Unix())
}

func (s *InMemorySession) Get(_ *http.Request, _ string) (map[any]any, error) {
	return s.Data, nil
}

func (s *InMemorySession) GetData(_ *http.Request, key ContextKey) (any, error) {
	if _, exists := s.Data[string(key)]; !exists {
		return nil, nil
	}

	if s.ExpiresAt < int(time.Now().Unix()) {
		return nil, errors.New("session expired")
	}

	return s.Data[string(key)], nil
}

func (s *InMemorySession) SetData(_ *http.Request, _ http.ResponseWriter, key ContextKey, value any) error {
	s.Data[string(key)] = value
	return nil
}
