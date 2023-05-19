package database

import (
	"log"
	"os"
	"testing"
)

var db *PostgreSQL

func TestMain(m *testing.M) {
	db = &PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}

	exitCode := m.Run()

	db.Close()

	os.Exit(exitCode)
}
