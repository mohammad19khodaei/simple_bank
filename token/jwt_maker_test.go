package token_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammad19khodaei/simple_bank/token"
	"github.com/mohammad19khodaei/simple_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	jwtMaker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := utils.RandomOwner()
	duration := time.Minute
	tokenString, err := jwtMaker.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := jwtMaker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now(), payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt.Time, time.Second)
}

func TestJWTMakerWithExpiredToken(t *testing.T) {
	jwtMaker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	tokenString, err := jwtMaker.GenerateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := jwtMaker.VerifyToken(tokenString)
	require.Error(t, err)
	require.True(t, errors.Is(err, jwt.ErrTokenExpired))
	require.Empty(t, payload)
}

func TestJWTMakerInvalidToken(t *testing.T) {
	jwtMaker, err := token.NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	payload, err := token.NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err = jwtMaker.VerifyToken(tokenString)
	require.Error(t, err)
	require.Empty(t, payload)
	require.True(t, errors.Is(err, token.ErrInvalidToken))
}
