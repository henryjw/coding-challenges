package main

import (
	"log"
	"rate-limiter/m/v2/server"
)

func main() {
	s := server.New()

	err := s.Run(8080)

	if err != nil {
		log.Fatalln("Error running the server: ", err)
	}
}
