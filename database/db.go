package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	DB *sql.DB
}

func (p *PostgreSQL) Connect(dbName string) error {
	dsn := fmt.Sprintf("user=postgres password=postgres dbname=%s sslmode=disable", dbName)
	var err error
	p.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := p.DB.Ping(); err != nil {
		p.DB.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (p *PostgreSQL) Close() error {
	return p.DB.Close()
}
