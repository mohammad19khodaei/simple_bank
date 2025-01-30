package db_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) db.Account {
	params := db.CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	params := db.UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, params.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := db.ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, accounts, 10)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
