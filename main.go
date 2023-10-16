package main

import (
	"database/sql"
	"log"

	"github.com/kellemNegasi/bank-system/api"
	db "github.com/kellemNegasi/bank-system/db/sqlc"
	"github.com/kellemNegasi/bank-system/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("unable to connect to the data base: %s", err)
	}

	store := db.NewStore(conn)
	server := api.New(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("couldn't start server: %v", err)
	}
}
