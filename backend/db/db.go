// this is for broad db function like opening a connection, closing a connection, and pinging.
// This may turn into a more advanced library inthe future, but for now, keep it simple stupid.

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Ping the database to check if the connection is valid
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
