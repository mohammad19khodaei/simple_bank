package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

var (
	testQueries *db.Queries
	testPool    *pgxpool.Pool
)

func TestMain(t *testing.M) {
	config, err := utils.LoadConfig("../..", "app")
	if err != nil {
		log.Fatal("could not load config", err)
	}

	ctx := context.Background()
	testPool, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = db.New(testPool)
	defer testPool.Close()
	existCode := t.Run()

	testQueries.DeleteAllTransfers(ctx)
	testQueries.DeleteAllEntries(ctx)
	testQueries.DeleteAllAccounts(ctx)

	os.Exit(existCode)
}
