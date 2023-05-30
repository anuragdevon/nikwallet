package services

import (
	"log"
	"nikwallet/config"
	"nikwallet/repository"
	"os"
	"testing"
)

var db *repository.PostgreSQL

func TestMain(m *testing.M) {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	db = &repository.PostgreSQL{}
	err = db.Connect(&c)
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}

	exitCode := m.Run()

	db.Close()

	os.Exit(exitCode)
}
