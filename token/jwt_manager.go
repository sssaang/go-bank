package token

import "fmt"

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