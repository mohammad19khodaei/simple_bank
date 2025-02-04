package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type Payload struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // Expiration time (exp)
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload, nil
}
