package data

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"password-manager/utils"
)

type Record struct {
	gorm.Model
	VaultId           uint
	Name              string
	UserName          string
	EncryptedPassword string
}

func CreateRecord(vaultId uint, recordName string, username string, password string) error {
	session := utils.GetSession()

	if session == nil {
		panic("Session should not be nil")
	}

	if recordExists(vaultId, recordName) {
		return errors.New(fmt.Sprintf("A record with the name '%v' already exists\n", recordName))
	}

	encryptedPassword, err := utils.Encrypt(password, session.VaultPasswordHash)

	if err != nil {
		panic(fmt.Sprintf("Error encrypting password: %v\n", err))
	}

	record := Record{
		VaultId:           vaultId,
		Name:              recordName,
		UserName:          username,
		EncryptedPassword: encryptedPassword,
	}
	db := GetDatabase()

	// TODO: execute this in a GoRoutine
	result := db.Create(&record)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetRecord(vaultId uint, recordName string) *Record {
	session := utils.GetSession()

	if session == nil {
		panic("Session should not be nil")
	}

	var record Record

	db := GetDatabase()

	db.Where(&Record{VaultId: vaultId, Name: recordName}).First(&record)

	if record.ID == 0 {
		return nil
	}

	return &record
}

func recordExists(vaultId uint, recordName string) bool {
	var record Record

	db := GetDatabase()

	db.Where(&Record{VaultId: vaultId, Name: recordName}).First(&record)

	return record.ID > 0
}
