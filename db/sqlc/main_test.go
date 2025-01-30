package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

// urlExample := "postgres://username:password@localhost:5432/database_name"
const (
	connectionString = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	testQueries *db.Queries
	testPool    *pgxpool.Pool
)

func TestMain(t *testing.M) {
	var err error
	ctx := context.Background()

	testPool, err = pgxpool.New(ctx, connectionString)
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
