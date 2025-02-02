package utils_test

import (
	"testing"

	"github.com/mohammad19khodaei/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := utils.RandomString(6)

	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isValid := utils.IsHashPasswordValid(hashedPassword, password)
	require.NoError(t, err)
	require.True(t, isValid)

	wrongPassword := utils.RandomString(6)
	isValid = utils.IsHashPasswordValid(wrongPassword, password)
	require.False(t, isValid)
}
