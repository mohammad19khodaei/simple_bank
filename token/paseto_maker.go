package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto    *paseto.V2
	secretKey []byte
}

func NewPasetoMaker(secretKey string) (Maker, error) {
	if len(secretKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid secret key size, must be exact %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:    &paseto.V2{},
		secretKey: []byte(secretKey),
	}

	return maker, nil
}

func (m *PasetoMaker) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.secretKey, payload, nil)
}

func (m *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	payload := &Payload{}

	err := m.paseto.Decrypt(tokenString, m.secretKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
