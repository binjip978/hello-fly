package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	s := newServer("")
	srv := httptest.NewServer(s.srv.Handler)
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
	s := newServer("")
	srv := httptest.NewServer(s.srv.Handler)
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

func TestSecret(t *testing.T) {
	s := newServer("")
	s.password = "password"
	s.user = "admin"
	srv := httptest.NewServer(s.srv.Handler)
	defer srv.Close()

	c := map[string]float64{
		"lat": 60.930432,
		"lon": 25.395363,
	}
	b, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.Post(srv.URL+"/secret", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, r.StatusCode)
	}

	req, err := http.NewRequest("POST", srv.URL+"/secret", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth("admin", "password")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
