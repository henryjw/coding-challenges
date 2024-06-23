package data

import (
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
	db := GetDatabase()

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
