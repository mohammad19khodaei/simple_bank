package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mohammad19khodaei/simple_bank/api"
	mockdb "github.com/mohammad19khodaei/simple_bank/db/mock"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransfer(t *testing.T) {
	fromAccount := createRandomAccount("USD")
	toAccount := createRandomAccount("USD")
	toEURAccount := createRandomAccount("EUR")

	testCases := []struct {
		name          string
		fromAccountID int32
		toAccountID   int32
		amount        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "from account id does not exists",
			fromAccountID: fromAccount.ID,
			toAccountID:   toAccount.ID,
			amount:        10,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "amount is greater than from account balance",
			fromAccountID: fromAccount.ID,
			toAccountID:   toAccount.ID,
			amount:        fromAccount.Balance + 10,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusPaymentRequired, recorder.Code)
			},
		},
		{
			name:          "to account id does not exists",
			fromAccountID: fromAccount.ID,
			toAccountID:   toAccount.ID,
			amount:        10,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), toAccount.ID).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "mismatch currency",
			fromAccountID: fromAccount.ID,
			toAccountID:   toEURAccount.ID,
			amount:        10,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), toEURAccount.ID).
					Times(1).
					Return(toEURAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "ok",
			fromAccountID: fromAccount.ID,
			toAccountID:   toAccount.ID,
			amount:        10,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), toAccount.ID).
					Times(1).
					Return(toAccount, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        10,
					}).
					Times(1).
					Return(createStubTransfer(fromAccount, toAccount, 10), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)

		server := api.NewServer(store)
		recorder := httptest.NewRecorder()
		requestBody := transferRequest{
			FromAccountID: tc.fromAccountID,
			ToAccountID:   tc.toAccountID,
			Amount:        tc.amount,
		}
		jsonData, err := json.Marshal(requestBody)
		fmt.Println(string(jsonData))
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(jsonData))
		request.Header.Set("Content-Type", "application/json")

		server.GetRouter().ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}
}

type transferRequest struct {
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID   int32 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TODO Need refactor
func createStubTransfer(fromAccount, toAccount db.Account, amount int) db.TransferTxResult {
	return db.TransferTxResult{
		Transfer:    db.Transfer{},
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		FromEntry:   db.Entry{},
		ToEntry:     db.Entry{},
	}
}
