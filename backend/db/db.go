// This is for broad db functions like opening a connection, closing a connection, and pinging.
// This may turn into a more advanced library in the future, but for now, keep it simple stupid.
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	var err error
	connectionString := "host=localhost port=5432 user=borea password=borea dbname=pg_borea sslmode=disable"

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
