package data

import (
	"fmt"
	"gorm.io/gorm"
	"password-manager/utils"
)

type Record struct {
	gorm.Model
	VaultId           uint
	Name              string
	EncryptedPassword string
}

func CreateRecord(vaultId uint, recordName string, recordPassword string) error {
	session := utils.GetSession()

	if session == nil {
		panic("Session should not be nil")
	}

	encryptedPassword, err := utils.Encrypt(recordPassword, session.VaultPasswordHash)

	if err != nil {
		panic(fmt.Sprintf("Error encrypting password: %v\n", err))
	}

	record := Record{Name: recordName, VaultId: vaultId, EncryptedPassword: encryptedPassword}
	db := GetDatabase()

	// TODO: execute this in a GoRoutine
	result := db.Create(&record)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
