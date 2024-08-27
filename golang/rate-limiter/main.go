package main

import (
	"log"
	rateLimiter "rate-limiter/m/v2/rate-limiter"
	"rate-limiter/m/v2/server"
)

func main() {
	limiter, err := rateLimiter.New(rateLimiter.Config{
		Algorithm:                   rateLimiter.RateLimitTokenBucket,
		MaxAllowedRequestsPerMinute: 3,
	})

	if err != nil {
		log.Fatalf("Unexpected error instantiating ratelimiter: %v\n", err)
	}

	s := server.New(limiter)

	err = s.Run(8080)

	if err != nil {
		log.Fatalln("Error running the server: ", err)
	}
}
