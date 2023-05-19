package handlers_test

import (
	"log"
	"nikwallet/database"
	"os"
	"testing"
)

var db *database.PostgreSQL

func TestMain(m *testing.M) {
	db = &database.PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}

	exitCode := m.Run()

	db.Close()

	os.Exit(exitCode)
}
