package db

import (
	"os"
	"testing"
)

const testDBName = "testdb"

func TestConnectToDBToCreateValidConnection(t *testing.T) {
	os.Setenv("DATABASE_NAME", testDBName)
	defer os.Unsetenv("DATABASE_NAME")

	db, err := ConnectToDB(testDBName)
	if err != nil {
		t.Fatalf("failed to create database instance: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}
