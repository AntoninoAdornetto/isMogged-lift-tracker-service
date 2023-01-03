package main

import (
	"database/sql"
	"log"

	"github.com/AntoninoAdornetto/lift_tracker/api"
	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config file", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server failed to start", err)
	}
}
