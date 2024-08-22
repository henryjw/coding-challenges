package server

import (
	"fmt"
	"net/http"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(portNumber int) error {
	http.HandleFunc("/unlimited", unlimitedHandler)
	http.HandleFunc("/limited", limitedHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil); err != nil {
		return nil
	}

	return nil
}

func unlimitedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All good"))
}

func limitedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:  add rate limiter
	w.WriteHeader(429)
}
