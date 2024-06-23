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

	vault := data.GetVault("test", "password123")

	if vault == nil {
		println("Vault not found")
		//data.CreateVault("test", "password123")
	} else {
		println("Vault found")
	}
}
