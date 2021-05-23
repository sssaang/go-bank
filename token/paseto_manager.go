package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoManager struct {
	paseto *paseto.V2
	symmetricKey []byte
}

func NewPasetoManager(symmetricKey string) (TokenManager, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: key must be exactly %d characaters", chacha20poly1305.KeySize)
	}

	manager := &PasetoManager{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return manager, nil
}

func (manager *PasetoManager) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return manager.paseto.Encrypt([]byte(manager.symmetricKey), payload, nil)
}
func (manager *PasetoManager) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := manager.paseto.Decrypt(token, []byte(manager.symmetricKey), payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}