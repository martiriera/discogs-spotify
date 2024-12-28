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

func (s *InMemorySession) Get(r *http.Request, sessionName string) (map[any]any, error) {
	return s.Data, nil
}

func (s *InMemorySession) GetData(r *http.Request, key string) (any, error) {
	if _, exists := s.Data[key]; !exists {
		return nil, nil
	}

	if s.ExpiresAt < int(time.Now().Unix()) {
		return nil, errors.New("session expired")
	}

	return s.Data[key], nil
}

func (s *InMemorySession) SetData(r *http.Request, w http.ResponseWriter, key string, value any) error {
	s.Data[key] = value
	return nil
}
