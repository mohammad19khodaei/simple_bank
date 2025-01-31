package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammad19khodaei/simple_bank/api"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(connPool)
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("could not start server", err)
	}
}
