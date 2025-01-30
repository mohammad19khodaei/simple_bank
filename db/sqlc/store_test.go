package db_test

import (
	"context"
	"fmt"
	"testing"

	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferTx(t *testing.T) {
	store := db.NewStore(testPool)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before", account1.Balance, account2.Balance)

	num := 5
	amount := int64(10)
	resultChan := make(chan result, num)
	for i := 0; i < num; i++ {
		go func() {
			transferResult, err := store.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			resultChan <- result{data: transferResult, err: err}
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < num; i++ {
		transferResult := <-resultChan
		require.NoError(t, transferResult.err)
		require.NotEmpty(t, transferResult.data)

		// check transfer
		transfer := transferResult.data.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := transferResult.data.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := transferResult.data.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := transferResult.data.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := transferResult.data.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">> tx ", i, fromAccount.Balance, toAccount.Balance)

		// check account's balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.Positive(t, diff1)

		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= num)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-amount*int64(num), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+amount*int64(num), updatedAccount2.Balance)
}

type result struct {
	err  error
	data db.TransferTxResult
}

func TestCreateTransferTxDeadLock(t *testing.T) {
	store := db.NewStore(testPool)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before", account1.Balance, account2.Balance)

	num := 10
	amount := int64(10)
	errorChan := make(chan error, num)
	for i := 0; i < num; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errorChan <- err
		}()
	}

	for i := 0; i < num; i++ {
		err := <-errorChan
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
