package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver     = "postgres"
	dbSourceName = "postgresql://root:secret@localhost:5432/basic_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSourceName)
	if err != nil {
		log.Fatalf("unable to connect to the data base: %s", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
