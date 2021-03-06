package token

import (
	"testing"
	"time"

	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestPasetoManager(t *testing.T) {
	manager, err := NewPasetoManager(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)
	
	username := util.RandomEmail()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := manager.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := manager.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	manager, err := NewPasetoManager(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)
	
	username := util.RandomEmail()
	token, err := manager.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := manager.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// TODO: add invalid paseto token test case

// func TestInvalidPasetoToken(t *testing.T) {
// 	manager, err := NewPasetoManager(util.RandomString(chacha20poly1305.KeySize))
// 	require.NoError(t, err)
	
// 	username := util.RandomEmail()
// 	duration := time.Minute

// 	token, err := manager.CreateToken(username, duration)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	payload, err := manager.VerifyToken("fasdfdsa.fadsfdsfds.fdsafsd")
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrInvalidToken.Error())
// 	require.Nil(t, payload)
// }