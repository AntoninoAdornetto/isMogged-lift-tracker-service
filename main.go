package main

import (
	"database/sql"
	"log"

	"github.com/AntoninoAdornetto/lift_tracker/api"
	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:tempPassword@localhost:5432/lift_tracker?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Server failed to start", err)
	}
}
