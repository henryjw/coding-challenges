package main

import (
	"database/sql"
	"dropbox/m/v2/services/auth"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./users.db")

	if err != nil {
		log.Fatalf("Error starting the database: %v\n", err)
	}

	defer db.Close()

	createUserTable(db)

	authService := auth.New(db)

	server := NewServer(authService)

	server.Run(8080)
}

func createUserTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"username" TEXT NOT NULL UNIQUE,
		"password" TEXT NOT NULL
	);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created users table")
	stmt.Exec()
}
