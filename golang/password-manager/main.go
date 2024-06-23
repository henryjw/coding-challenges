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

	createVaultErr := data.CreateVault("test", "password123")

	if createVaultErr != nil {
		fmt.Printf("Error creating vault: %v\n", createVaultErr)
	}

	if vault == nil {
		println("Vault not found")
	} else {
		println("Vault found")
	}
}
