package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type server struct {
	srv      *http.Server
	user     string
	password string
}

func newServer(address string) *server {
	s := server{
		srv: &http.Server{
			Addr: address,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/hello", s.hello)
	mux.HandleFunc("/fly", s.fly)
	mux.HandleFunc("/secret", s.secret)
	mux.HandleFunc("/", s.hello)
	s.srv.Handler = mux

	return &s
}

func (s *server) healthz(w http.ResponseWriter, r *http.Request) {
	log.Println("/healthz")
	w.WriteHeader(http.StatusOK)
}

func (s *server) hello(w http.ResponseWriter, r *http.Request) {
	log.Println("serving hello request")
	v := rand.Float64()
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("hostname error: %v", err)
	}

	rsp := map[string]string{
		"info":     "simple http service",
		"number":   fmt.Sprintf("%.4f", v),
		"author":   "me",
		"hostname": hostname,
	}

	b, err := json.Marshal(&rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (s *server) fly(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><h1>Hello Fly v2</h1></html>"))
}

func (s *server) secret(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, password, ok := r.BasicAuth()
	// THIS IS SECURE
	if !ok || username != s.user || password != s.password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	coords := make(map[string]float64)
	err = json.Unmarshal(b, &coords)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// echo coordinates back
	b, err = json.Marshal(&coords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func main() {
	// SANE DEFAULTS
	user := os.Getenv("FLY_USER")
	if user == "" {
		user = "user"
	}
	password := os.Getenv("FLY_PASSWORD")
	if password == "" {
		password = "password"
	}

	s := newServer(":8080")
	s.user = user
	s.password = password
	log.Printf("starting web-server on :8080")
	log.Fatal(s.srv.ListenAndServe())
}
