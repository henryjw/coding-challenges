package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	rateLimiter "rate-limiter/m/v2/rate-limiter"
)

type Server struct {
	limiter rateLimiter.IRateLimiter
}

func New(limiter rateLimiter.IRateLimiter) *Server {
	return &Server{
		limiter: limiter,
	}
}

func (s *Server) Run(portNumber int) error {
	http.HandleFunc("/unlimited", unlimitedHandler)
	http.HandleFunc("/limited", s.limitedHandler)

	log.Printf("HTTP server running on port %d\n", portNumber)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil); err != nil {
		return nil
	}

	return nil
}

func unlimitedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unlimited! Let's Go!"))
}

func (s *Server) limitedHandler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, ipParseErr := net.SplitHostPort(r.RemoteAddr)

	if ipParseErr != nil {
		log.Printf("Unexpected error parsing client IP address: %v\n", ipParseErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error"))

		return
	}

	requestInfo := rateLimiter.RequestInfo{
		IPAddress: ipAddress,
		// TODO: remove query params
		Endpoint: r.RequestURI,
	}
	if allowed, err := s.limiter.AllowRequest(requestInfo); err == nil {
		if allowed {
			w.Write([]byte("Limited, don't over use me!"))
			return
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit reached!"))
			return
		}
	} else {
		log.Printf("Unexpected error: %v\n", err)
		w.Write([]byte("An error occurred"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
