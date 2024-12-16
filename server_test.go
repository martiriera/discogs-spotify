package main

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
}

func newPostListRequest(listUrl string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/list", nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
