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

type VaultManager struct {
	db *gorm.DB
}

func NewVaultManager(database *gorm.DB) *VaultManager {
	return &VaultManager{
		db: database,
	}
}

func (receiver *VaultManager) CreateVault(name string, password string) error {
	passwordHash := utils.Hash(password, name)
	vault := Vault{Name: name, PasswordHash: passwordHash}

	if receiver.vaultExists(name) {
		return errors.New(fmt.Sprintf("A vault with the name '%s' already exists", name))
	}

	// TODO: execute this in a GoRoutine
	result := receiver.db.Create(&vault)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (receiver *VaultManager) GetVault(name string, password string) *Vault {
	passwordHash := utils.Hash(password, name)
	var vault Vault

	receiver.db.Where(&Vault{Name: name, PasswordHash: passwordHash}).First(&vault)
	if vault.ID == 0 {
		return nil
	}

	return &vault
}

func (receiver *VaultManager) vaultExists(name string) bool {
	var vault Vault

	receiver.db.Where(&Vault{Name: name}).First(&vault)

	return vault.ID > 0
}
