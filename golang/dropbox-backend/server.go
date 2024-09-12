package main

import (
	"dropbox/m/v2/services/auth"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	authService *auth.AuthService
}

type ErrorResponse struct {
	Err string `json:"err"`
}

type loginResponse struct {
	Token string
}

func NewServer(authService *auth.AuthService) *Server {
	return &Server{
		authService: authService,
	}
}

func (s *Server) Run(portNumber int) error {
	http.HandleFunc("/login", s.loginHandler)

	log.Printf("HTTP server running on port %d\n", portNumber)

	err := http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)

	return err
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user auth.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.authService.Login(user)

	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, auth.InvalidLoginError) {
		res, _ := json.Marshal(ErrorResponse{
			Err: "Invalid username or password",
		})
		http.Error(w, string(res), http.StatusBadRequest)
		return
	}

	if err != nil {
		res, _ := json.Marshal(ErrorResponse{
			Err: fmt.Sprintf("Unexpected error: %s", err.Error()),
		})

		http.Error(w, string(res), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
