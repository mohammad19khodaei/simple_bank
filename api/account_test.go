package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mohammad19khodaei/simple_bank/api"
	"github.com/mohammad19khodaei/simple_bank/api/middlewares"
	mockdb "github.com/mohammad19khodaei/simple_bank/db/mock"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	account := createRandomAccount()

	testCases := []struct {
		name          string
		accountID     int32
		setAuthHeader func(maker token.Maker, req *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setAuthHeader: func(tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken(account.Owner, config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				gotAccount := db.Account{}
				json.Unmarshal(body, &gotAccount)
				require.Equal(t, account, gotAccount)
			},
		},
		{
			name:      "not found",
			accountID: account.ID,
			setAuthHeader: func(tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken(account.Owner, config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "internal error",
			accountID: account.ID,
			setAuthHeader: func(tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken(account.Owner, config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "invalid id",
			accountID: 0,
			setAuthHeader: func(tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken(account.Owner, config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	server, err := api.NewServer(config, store)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request := httptest.NewRequest(http.MethodGet, url, nil)

			tc.setAuthHeader(tokenMaker, request)
			server.Router().ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
