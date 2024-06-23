package main

import (
	"fmt"
	"password-manager/data"
)

func main() {
	dbErr := data.InitDatabase()

	if dbErr != nil {
		panic(fmt.Sprintf("Error initializing database: %v\n", dbErr))
	}

	data.CreateVault("test", "password123")
}
