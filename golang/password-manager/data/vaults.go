package data

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"password-manager/utils"
)

type Vault struct {
	gorm.Model
	Name         string
	PasswordHash string
}

func CreateVault(name string, password string) error {
	passwordHash := utils.Hash(password, name)
	vault := Vault{Name: name, PasswordHash: passwordHash}

	// TODO: consider creating a struct to maintain common state (like the database pointer and the session) so it can be injected
	// instead of having to retrieve it in every method
	db := GetDatabase()

	if vaultExists(name) {
		return errors.New(fmt.Sprintf("A vault with the name '%s' already exists", name))
	}

	// TODO: execute this in a GoRoutine
	result := db.Create(&vault)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetVault(name string, password string) *Vault {
	db := GetDatabase()
	passwordHash := utils.Hash(password, name)
	var vault Vault

	db.Where(&Vault{Name: name, PasswordHash: passwordHash}).First(&vault)
	if vault.ID == 0 {
		return nil
	}

	return &vault
}

func vaultExists(name string) bool {
	db := GetDatabase()

	var vault Vault

	db.Where(&Vault{Name: name}).First(&vault)

	return vault.ID > 0
}
