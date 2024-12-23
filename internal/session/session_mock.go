package session

import "net/http"

type InMemorySession struct {
	Data map[interface{}]interface{}
}

func NewInMemorySession() *InMemorySession {
	return &InMemorySession{
		Data: make(map[interface{}]interface{}),
	}
}

func (s *InMemorySession) Init() {}

func (s *InMemorySession) Get(r *http.Request, sessionName string) (*SessionData, error) {
	return &SessionData{Values: s.Data}, nil
}

func (s *InMemorySession) GetData(r *http.Request, key string) (interface{}, error) {
	return s.Data[key], nil
}

func (s *InMemorySession) SetData(r *http.Request, w http.ResponseWriter, key string, value interface{}) error {
	s.Data[key] = value
	return nil
}
