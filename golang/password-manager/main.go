package main

import (
	"fmt"
	"password-manager/data"
	"password-manager/utils"
)

func main() {
	db, dbErr := data.InitDatabase()

	if dbErr != nil {
		panic(fmt.Sprintf("Error initializing database: %v\n", dbErr))
	}

	vaultManager := data.NewVaultManager(db)
	recordManager := data.NewRecordManager(db)

	vault := vaultManager.GetVault("test", "password123")

	createVaultErr := vaultManager.CreateVault("test", "password123")

	if createVaultErr != nil {
		fmt.Printf("Error creating vault: %v\n", createVaultErr)
	}

	if vault == nil {
		panic("Vault not found")
	} else {
		println("Vault found")
	}

	session := utils.CreateSession(vault.ID, vault.PasswordHash)

	println("Vault password hash", session.VaultPasswordHash)

	recordName := "GoogleDrive"
	createRecordError := recordManager.CreateRecord(vault.ID, recordName, "henry", "secure_password123")

	if createRecordError != nil {
		fmt.Printf("Error creating record: %v\n", createRecordError)
	}

	record := recordManager.GetRecord(vault.ID, recordName)

	if record == nil {
		fmt.Printf("No record named '%v' found\n", recordName)
		return
	}

	fmt.Printf("Record: %v\n", record)
	password, decryptErr := utils.Decrypt(record.EncryptedPassword, session.VaultPasswordHash)

	if decryptErr != nil {
		fmt.Printf("Error decrypting password: %v\n", decryptErr)
	} else {
		fmt.Printf("Decrypted password: %s\n", password)
	}

}
