package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store is a wrapper around sql.DB that provides a set of methods
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to another
// It creates a new transfer record, updates the balance of two accounts,
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txKey := ctx.Value(txKey)

		fmt.Println(txKey, "create transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: int64(arg.FromAccountID),
			ToAccountID:   int64(arg.ToAccountID),
			Amount:        int64(arg.Amount),
		})
		if err != nil {
			return err
		}

		fmt.Println(txKey, "create from entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: int64(arg.FromAccountID),
			Amount:    -int64(arg.Amount),
		})

		if err != nil {
			return err
		}

		fmt.Println(txKey, "create to entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: int64(arg.ToAccountID),
			Amount:    int64(arg.Amount),
		})

		if err != nil {
			return err
		}


		fmt.Println(txKey, "get account 1")
		account1, err := q.GetAccountForUpdate(ctx, int64(arg.FromAccountID))
		if err != nil {
			return err
		}

		fmt.Println(txKey, "update account 1")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      int64(arg.FromAccountID),
			Balance: account1.Balance - int64(arg.Amount),
		})

		if err != nil {
			return err
		}

		fmt.Println(txKey, "get account 2")
		account2, err := q.GetAccountForUpdate(ctx, int64(arg.ToAccountID))
		if err != nil {
			return err
		}

		fmt.Println(txKey, "update account 2")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      int64(arg.ToAccountID),
			Balance: account2.Balance + int64(arg.Amount),
		})

		if err != nil {
			return err
		}



		return nil
	})
	return result, err
}
