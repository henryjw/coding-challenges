package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var logger = log.Default()

const databaseFileName = "data.sqlite"

var db *gorm.DB

func InitDatabase() error {
	var err error

	logger.Println("Initializing database...")
	db, err = gorm.Open(sqlite.Open(databaseFileName), &gorm.Config{})

	if err != nil {
		logger.Printf("Error initializing databasae: %v\n", err)
		return err
	}

	logger.Println("Database initialized")

	return nil
}

func GetDatabase() *gorm.DB {
	if db == nil {
		panic("Database has not been initialized")
	}

	return db
}
