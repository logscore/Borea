package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initializeDatabase() {
	// Open a database connection
	db, err := sql.Open("sqlite3", "../../db.sqlite")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	fmt.Println("Connected to the SQLite database.")

	// Create the 'sessions' table
	_, err = db.Exec(`
        CREATE TABLE sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sessionId TEXT,
            lastActivityTime TEXT,
            userId TEXT,
            userPath TEXT,
            sessionDuration INTEGER,
            userAgent TEXT,
            referrer TEXT,
            token TEXT,
            startTime TEXT,
            screenResolution TEXT,
            language TEXT,
            FOREIGN KEY (userId) REFERENCES users(userId)
        );`)
	if err != nil {
		log.Fatal("Error creating 'sessions' table:", err)
	}
	fmt.Println("'sessions' table created successfully.")

	// Create the 'admin_users' table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS admin_users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            Username TEXT NOT NULL,
            PasswordHash TEXT NOT NULL
        );`)
	if err != nil {
		log.Fatal("Error creating 'admin_users' table:", err)
	}
	fmt.Println("'admin_users' table created successfully.")

	// Create the 'unique_users' table
	_, err = db.Exec(`
            CREATE TABLE IF NOT EXISTS unique_users (
                userId TEXT PRIMARY KEY UNIQUE,
                lastActivityTime TEXT
            );
        `)
	if err != nil {
		log.Fatal("Error creating 'unique_users' table:", err)
	}
	fmt.Println("'unique_users' table created successfully.")

	fmt.Println("Database initialization completed.")
}

func main() {
	initializeDatabase()
}
