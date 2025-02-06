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
	"github.com/mohammad19khodaei/simple_bank/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		params        createUserParams
		buildStubs    func(store *mockdb.MockStore, params createUserParams)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, params createUserParams)
	}{
		{
			name:   "without any params",
			params: createUserParams{},
			buildStubs: func(store *mockdb.MockStore, _ createUserParams) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ createUserParams) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// TODO check response body when you refactor the error response messages
			},
		},
		{
			name: "without username",
			params: createUserParams{
				FullName: utils.RandomOwner(),
				Email:    utils.RandomEmail(),
				Password: utils.RandomString(6),
			},
			buildStubs: func(store *mockdb.MockStore, _ createUserParams) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ createUserParams) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// TODO add more test for checking other params
		{
			name: "with correct params",
			params: createUserParams{
				Username: utils.RandomOwner(),
				FullName: utils.RandomOwner(),
				Email:    utils.RandomEmail(),
				Password: utils.RandomString(6),
			},
			buildStubs: func(store *mockdb.MockStore, params createUserParams) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(_ context.Context, arg db.CreateUserParams) (db.User, error) {
						require.True(t, utils.IsHashPasswordValid(arg.HashedPassword, params.Password))

						return db.User{
							Username:          arg.Username,
							HashedPassword:    arg.HashedPassword,
							FullName:          arg.FullName,
							Email:             arg.Email,
							PasswordChangedAt: pgtype.Timestamptz{}, // Mock time value if necessary
							CreatedAt:         pgtype.Timestamptz{}, // Mock time value if necessary
						}, nil
					})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, params createUserParams) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var resp api.UserResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, params.Username, resp.Username)
				require.Equal(t, params.FullName, resp.FullName)
				require.Equal(t, params.Email, resp.Email)
			},
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store, tc.params)
		server, err := api.NewServer(config, store)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()
		jsonData, err := json.Marshal(tc.params)
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonData))

		server.Router().ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder, tc.params)
	}

}

type createUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func TestLogin(t *testing.T) {
	password := utils.RandomString(6)
	user := createRandomUser(password)
	testCases := []struct {
		name          string
		params        loginParams
		buildStubs    func(store *mockdb.MockStore, params loginParams)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, params loginParams)
	}{
		{
			name: "username does not exists",
			params: loginParams{
				Username: utils.RandomOwner(),
				Password: utils.RandomString(6),
			},
			buildStubs: func(store *mockdb.MockStore, params loginParams) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(params.Username)).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ loginParams) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "password is wrong",
			params: loginParams{
				Username: user.Username,
				Password: utils.RandomString(6),
			},
			buildStubs: func(store *mockdb.MockStore, params loginParams) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(params.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, _ loginParams) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ok",
			params: loginParams{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, params loginParams) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(params.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, params loginParams) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var resp api.LoginResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.NotEmpty(t, resp.AccessToken)
				require.Equal(t, params.Username, resp.User.Username)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	server, err := api.NewServer(config, store)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store, tc.params)

			recorder := httptest.NewRecorder()
			jsonData, err := json.Marshal(tc.params)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(jsonData))

			server.Router().ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, tc.params)
		})
	}
}

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
