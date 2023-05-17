package database

import (
	"testing"
)

func TestPostgreSQL(t *testing.T) {
	t.Run("Connect to database successfully with valid dbname", func(t *testing.T) {
		db := &PostgreSQL{}
		err := db.Connect("testdb")
		if err != nil {
			t.Errorf("error connecting to database: %s", err)
		}
		defer db.Close()

		if db.DB == nil {
			t.Error("failed to connect to database")
		}
	})

	t.Run("Close connection to db successfully", func(t *testing.T) {
		db := &PostgreSQL{}
		err := db.Connect("testdb")
		if err != nil {
			t.Errorf("error connecting to database: %s", err)
		}

		err = db.Close()
		if err != nil {
			t.Errorf("error closing database connection: %s", err)
		}
	})

	t.Run("Connect to database to return error with invalid dbName", func(t *testing.T) {
		db := &PostgreSQL{}
		err := db.Connect("invalid_db")
		if err == nil {
			t.Error("expected an error connecting to invalid database, but got nil")
		}
	})
}
