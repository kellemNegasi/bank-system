package main

import (
	"database/sql"
	"log"

	"github.com/kellemNegasi/bank-system/api"
	db "github.com/kellemNegasi/bank-system/db/sqlc"
	_ "github.com/lib/pq"
)

var address string = "0.0.0.0:8080"

const (
	dbDriver     = "postgres"
	dbSourceName = "postgresql://root:secret@localhost:5432/basic_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSourceName)
	if err != nil {
		log.Fatalf("unable to connect to the data base: %s", err)
	}

	store := db.NewStore(conn)
	server := api.New(store)
	err = server.Start(address)
	if err != nil {
		log.Fatalf("couldn't start server: %v", err)
	}
}
