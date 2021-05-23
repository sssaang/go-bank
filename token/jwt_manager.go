package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const MinSecretLength = 32

type JWTManager struct {
	Secret string
}

func NewJWTManager(secret string) (TokenManager, error) {
	if len(secret) < MinSecretLength {
		return nil, fmt.Errorf("invalid key size: must be at least %d character", MinSecretLength)
	}

	return &JWTManager {
		Secret: secret,
	}, nil
}

func (manager JWTManager) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, payload)
	return jwtToken.SignedString([]byte(manager.Secret))
}

func (manager JWTManager) VerifyToken(token string) (*Payload, error) {
	return &Payload{}, nil
}

