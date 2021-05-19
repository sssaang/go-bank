package db

import (
	"context"
	"testing"

	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
)
func createRandomUser(t *testing.T) User {
	username := util.RandomEmail()
	arg := CreateUserParams{
		Username: username,
		HashedPassword: "pw",
		FullName: util.RandomOwner(),
		Email: username,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	userCreated := createRandomUser(t)
	userFetched, err := testQueries.GetUser(context.Background(), userCreated.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userFetched)

	require.Equal(t, userCreated.Username, userFetched.Username)
	require.Equal(t, userCreated.HashedPassword, userFetched.HashedPassword)
	require.Equal(t, userCreated.FullName, userFetched.FullName)
	require.Equal(t, userCreated.Email, userFetched.Email)
	require.WithinDuration(t, userCreated.PasswordChangedAt, userFetched.PasswordChangedAt, 0)
	require.WithinDuration(t, userCreated.CreatedAt, userFetched.CreatedAt, 0)
}