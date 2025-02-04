package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		params        transferRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, param transferRequest)
	}{
		{
			name: "from account id does not exists",
			params: transferRequest{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ transferRequest) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "amount is greater than from account balance",
			params: transferRequest{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        fromAccount.Balance + 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccount.ID).
					Times(1).
					Return(fromAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ transferRequest) {
				require.Equal(t, http.StatusPaymentRequired, recorder.Code)
			},
		},
		{
			name: "to account id does not exists",
			params: transferRequest{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        10,
			},
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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ transferRequest) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "mismatch currency",
			params: transferRequest{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toEURAccount.ID,
				Amount:        10,
			},
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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ transferRequest) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ok",
			params: transferRequest{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        10,
			},
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
					DoAndReturn(func(_ context.Context, arg db.TransferTxParams) (db.TransferTxResult, error) {
						return db.TransferTxResult{
							Transfer: db.Transfer{
								ID:            1,
								FromAccountID: arg.FromAccountID,
								ToAccountID:   arg.ToAccountID,
								Amount:        arg.Amount,
								CreatedAt:     pgtype.Timestamptz{},
							},
							FromAccount: fromAccount,
							ToAccount:   toAccount,
							FromEntry: db.Entry{
								ID:        1,
								AccountID: arg.FromAccountID,
								Amount:    arg.Amount,
								CreatedAt: pgtype.Timestamptz{},
							},
							ToEntry: db.Entry{
								ID:        2,
								AccountID: arg.ToAccountID,
								Amount:    arg.Amount,
								CreatedAt: pgtype.Timestamptz{},
							},
						}, nil
					})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, param transferRequest) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var resp db.TransferTxResult
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, resp.Transfer.FromAccountID, param.FromAccountID)
				require.Equal(t, resp.Transfer.ToAccountID, param.ToAccountID)
				require.Equal(t, resp.Transfer.Amount, param.Amount)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := api.NewServer(config, store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			jsonData, err := json.Marshal(tc.params)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(jsonData))
			request.Header.Set("Content-Type", "application/json")

			server.Router().ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, tc.params)
		})
	}
}

type transferRequest struct {
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID   int32 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
