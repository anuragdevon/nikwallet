package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"nikwallet/pkg/db"
	"os"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := db.ConnectToDB("testdb"); err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	testDB = db.DB

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func JsonEqual(a, b []byte) bool {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false
	}
	return bytes.Equal(a, b)
}
