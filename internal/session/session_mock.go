package session

import "net/http"

type InMemorySession struct {
	Data map[any]any
}

func NewInMemorySession() *InMemorySession {
	return &InMemorySession{
		Data: make(map[any]any),
	}
}

func (s *InMemorySession) Init() {}

func (s *InMemorySession) Get(r *http.Request, sessionName string) (map[any]any, error) {
	return s.Data, nil
}

func (s *InMemorySession) GetData(r *http.Request, key string) (any, error) {
	return s.Data[key], nil
}

func (s *InMemorySession) SetData(r *http.Request, w http.ResponseWriter, key string, value any) error {
	s.Data[key] = value
	return nil
}
