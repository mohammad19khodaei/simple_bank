package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammad19khodaei/simple_bank/api"
)

const (
	connectionString = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address          = "0.0.0.0:8080"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(connPool)
	if err := server.Start(address); err != nil {
		log.Fatal("could not start server", err)
	}
}
