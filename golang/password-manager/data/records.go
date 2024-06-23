package data

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	Vault             Vault
	Name              string
	EncryptedPassword string
}

func CreateRecord(vaultId uint, recordName string, recordPassword string) error {

}
