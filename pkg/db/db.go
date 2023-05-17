package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectToDB(dbName string) error {
	dsn := fmt.Sprintf("user=postgres password=postgres dbname=%s sslmode=disable", dbName)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		DB.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
