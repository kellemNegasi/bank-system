package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kellemNegasi/bank-system/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../../")
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("unable to connect to the data base: %s", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
