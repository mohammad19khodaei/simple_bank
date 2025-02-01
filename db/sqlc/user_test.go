package db_test

import (
	"context"
	"testing"

	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) db.User {
	params := db.CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: "secret",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}
