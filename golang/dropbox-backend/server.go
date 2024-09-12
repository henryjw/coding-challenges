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
	http.HandleFunc("/signup", s.signUpHandler)

	log.Printf("HTTP server running on port %d\n", portNumber)

	err := http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)

	return err
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user auth.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		sendJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateUser(user)

	if err != nil {
		sendJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.authService.Login(user)

	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, auth.InvalidLoginError) {
		sendJsonError(w, "Invalid username or password", http.StatusInternalServerError)
		return
	}

	if err != nil {
		sendJsonError(w, fmt.Sprintf("Unexpected error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user auth.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		sendJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateUser(user)

	if err != nil {
		sendJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.authService.SignUp(user)

	if err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, auth.UserAlreadyExistsError) {
			statusCode = 400
		}

		sendJsonError(w, err.Error(), statusCode)
	}

	w.WriteHeader(http.StatusCreated)
}

func sendJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	res, _ := json.Marshal(ErrorResponse{
		Err: message,
	})

	http.Error(w, string(res), statusCode)
}

func validateUser(user auth.User) error {
	if user.Password == "" || user.Username == "" {
		return errors.New("`username` and `password` fields are required")
	}

	return nil
}
