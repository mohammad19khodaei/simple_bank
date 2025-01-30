package db_test

import (
	"context"
	"testing"

	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	params := db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)
}
