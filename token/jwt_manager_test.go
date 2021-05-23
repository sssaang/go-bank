package token

import (
	"testing"
	"time"

	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	manager, err := NewJWTManager(util.RandomString(32))
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