package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	srv := httptest.NewServer(mux())
	defer srv.Close()

	r, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, r.StatusCode)
	}
}

func TestHello(t *testing.T) {
	srv := httptest.NewServer(mux())
	defer srv.Close()

	r, err := http.Get(srv.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, r.StatusCode)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	js := make(map[string]any)
	err = json.Unmarshal(b, &js)
	if err != nil {
		t.Fatal(err)
	}

	author := js["author"].(string)
	if author != "me" {
		t.Errorf("expected me, go %s", author)
	}
}
