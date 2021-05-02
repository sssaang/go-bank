package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
	return transfer
}

func TestDeleteTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), fromAccount.ID)
	defer testQueries.DeleteAccount(context.Background(), toAccount.ID)

	transferCreated := CreateRandomTransfer(t, fromAccount, toAccount)

	err := testQueries.DeleteTransfer(context.Background(), transferCreated.ID)
	require.NoError(t, err)

	transferDeleted, err := testQueries.GetTransfer(context.Background(), transferCreated.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transferDeleted)
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), fromAccount.ID)
	defer testQueries.DeleteAccount(context.Background(), toAccount.ID)
	
	transferCreated := CreateRandomTransfer(t, fromAccount, toAccount)
	defer testQueries.DeleteTransfer(context.Background(), transferCreated.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), fromAccount.ID)
	defer testQueries.DeleteAccount(context.Background(), toAccount.ID)
	
	transferCreated := CreateRandomTransfer(t, fromAccount, toAccount)
	defer testQueries.DeleteTransfer(context.Background(), transferCreated.ID)
	
	transferFetched, err := testQueries.GetTransfer(context.Background(), transferCreated.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferFetched)

	require.Equal(t, transferCreated.ID, transferFetched.ID)
	require.Equal(t, transferCreated.FromAccountID, transferFetched.FromAccountID)
	require.Equal(t, transferCreated.ToAccountID, transferFetched.ToAccountID)
	require.Equal(t, transferCreated.Amount, transferFetched.Amount)
	require.WithinDuration(t, transferCreated.CreatedAt, transferFetched.CreatedAt, 0)
}

func TestListTransfers(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), fromAccount.ID)
	defer testQueries.DeleteAccount(context.Background(), toAccount.ID)

	for i := 0; i < 10; i++ {
		transfer := CreateRandomTransfer(t, fromAccount, toAccount)
		defer testQueries.DeleteTransfer(context.Background(), transfer.ID)
	}


	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Limit: 5,
		Offset: 5,
	}

	transferList, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transferList, 5)

	for _, transfer := range transferList {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
	}
}
