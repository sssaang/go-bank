package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sssaang/simplebank/api"
	db "github.com/sssaang/simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:test@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:1234"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}  