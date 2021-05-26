package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestDeleteAccount(t *testing.T) {
	accountCreated := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), accountCreated.ID)
	require.NoError(t, err)

	accountDeleted, err := testQueries.GetAccount(context.Background(), accountCreated.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountDeleted)
}

func TestCreateAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), account.ID)
}

func TestGetAccount(t *testing.T) {
	accountCreated := CreateRandomAccount(t)
	accountFetched, err := testQueries.GetAccount(context.Background(), accountCreated.ID)
	defer testQueries.DeleteAccount(context.Background(), accountCreated.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFetched)

	require.Equal(t, accountCreated.ID, accountFetched.ID)
	require.Equal(t, accountCreated.Owner, accountFetched.Owner)
	require.Equal(t, accountCreated.Balance, accountFetched.Balance)
	require.Equal(t, accountCreated.Currency, accountFetched.Currency)
	require.WithinDuration(t, accountCreated.CreatedAt, accountFetched.CreatedAt, 0)
}

func TestListAccounts(t *testing.T) {
	var testAccount Account

	for i := 0; i < 10; i++ {
		testAccount = CreateRandomAccount(t)
		defer testQueries.DeleteAccount(context.Background(), testAccount.ID)
	}

	arg := ListAccountsParams{
		Owner: testAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	accountList, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountList)

	for _, account := range accountList {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, testAccount.Owner)
	}
}

func TestUpdateAccount(t *testing.T) {
	accountCreated := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), accountCreated.ID)
	
	arg := UpdateAccountParams{
		ID: accountCreated.ID,
		Balance: util.RandomMoney(),
	}

	accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)	

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, accountCreated.ID, accountUpdated.ID)
	require.Equal(t, accountCreated.Owner, accountUpdated.Owner)
	require.Equal(t, accountUpdated.Balance, arg.Balance)
	require.Equal(t, accountCreated.Currency, accountUpdated.Currency)

}
