package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var logger = log.Default()

const databaseFileName = "data.sqlite"

func InitDatabase() (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	logger.Println("Initializing database...")
	db, err = gorm.Open(sqlite.Open(databaseFileName), &gorm.Config{})

	if err != nil {
		logger.Printf("Error initializing databasae: %v\n", err)
		return nil, err
	}

	migrateError := db.AutoMigrate(&Vault{}, &Record{})

	if migrateError != nil {
		logger.Printf("Error running auto-migration: %v\n", migrateError)
		return nil, migrateError
	}

	logger.Println("Database initialized")

	return db, nil
}
