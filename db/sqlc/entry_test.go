package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, arg.Amount)
	return entry
}

func TestDeleteEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), account.ID)
	
	entryCreated := CreateRandomEntry(t, account)

	err := testQueries.DeleteEntry(context.Background(), entryCreated.ID)
	require.NoError(t, err)

	entryDeleted, err := testQueries.GetEntry(context.Background(), entryCreated.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entryDeleted)
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), account.ID)

	entry := CreateRandomEntry(t, account)
	defer testQueries.DeleteEntry(context.Background(), entry.ID)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), account.ID)

	entryCreated := CreateRandomEntry(t, account)
	defer testQueries.DeleteEntry(context.Background(), entryCreated.ID)

	entryFetched, err := testQueries.GetEntry(context.Background(), entryCreated.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryFetched)

	require.Equal(t, entryCreated.ID, entryFetched.ID)
	require.Equal(t, entryCreated.AccountID, entryFetched.AccountID)
	require.Equal(t, entryCreated.Amount, entryFetched.Amount)
	require.WithinDuration(t, entryCreated.CreatedAt, entryFetched.CreatedAt, 0)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	defer testQueries.DeleteAccount(context.Background(), account.ID)

	for i := 0; i < 10; i++ {
		entry := CreateRandomEntry(t, account)
		defer testQueries.DeleteEntry(context.Background(), entry.ID)
	}


	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}

	entryList, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entryList, 5)

	for _, entry := range entryList {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account.ID)
	}
}
