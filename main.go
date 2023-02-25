package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type server struct {
	srv *http.Server
}

func newServer(address string, mux http.Handler) server {
	return server{
		srv: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Println("/healthz")
	w.WriteHeader(http.StatusOK)
}

func hello(w http.ResponseWriter, r *http.Request) {
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

func fly(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><h1>Hello Fly v2</h1></html>"))
}

func musk(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("nano-manager"))
}

func mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/fly", fly)
	mux.HandleFunc("/", hello)
	return mux
}

func main() {
	s := newServer(":8080", mux())
	log.Printf("starting web-server on :8080")
	log.Fatal(s.srv.ListenAndServe())
}
