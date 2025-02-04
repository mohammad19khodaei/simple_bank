package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammad19khodaei/simple_bank/api"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".", "app")
	if err != nil {
		log.Fatal("could not load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	server, err := api.NewServer(config, db.NewStore(connPool))
	if err != nil {
		log.Fatal("could not create start", err)
	}

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("could not start server", err)
	}
}
