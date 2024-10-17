// This is for broad db functions like opening a connection, closing a connection, and pinging.
// This may turn into a more advanced library in the future, but for now, keep it simple stupid.
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	PG_HOST := os.Getenv("PG_HOST")
	PG_PORT := os.Getenv("PG_PORT")
	PG_USER := os.Getenv("PG_USER")
	PG_PSWD := os.Getenv("PG_PSWD")
	DB_NAME := os.Getenv("DB_NAME")

	var err error
	connectionString := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, PG_HOST, PG_PORT, PG_USER, PG_PSWD, DB_NAME)

	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Ping the database to check if the connection is valid
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
