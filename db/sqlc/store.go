package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store {
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error)  error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err:%v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO Update Money w/o being trapped in deadlocks
		result.FromAccount, result.ToAccount, err = TransferMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)
		
		if err != nil {
			return err
		}
		
		return nil
	})

	return result, err
}

func TransferMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (fromAccount Account, toAccount Account, err error) {

	if fromAccountID < toAccountID {
		fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams {
			ID: fromAccountID,
			Amount: -amount,
		})

		if err != nil {
			return 
		}

		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams {
			ID: toAccountID,
			Amount: amount,
		}) 

	} else {
		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams {
			ID: toAccountID,
			Amount: amount,
		})

		if err != nil {
			return 
		}

		fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams {
			ID: fromAccountID,
			Amount: -amount,
		})
  
	}	

	return
}