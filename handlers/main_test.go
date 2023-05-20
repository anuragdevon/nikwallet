package handlers

import (
	"log"
	"nikwallet/repository"
	"os"
	"testing"
)

var db *repository.PostgreSQL

func TestMain(m *testing.M) {
	db = &repository.PostgreSQL{}
	err := db.Connect("testdb")
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}

	exitCode := m.Run()

	db.Close()

	os.Exit(exitCode)
}
