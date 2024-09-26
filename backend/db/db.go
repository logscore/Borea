package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "../db.sqlite")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Optionally, ping the database to check if the connection is valid
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection successfully established")
	return nil
}
