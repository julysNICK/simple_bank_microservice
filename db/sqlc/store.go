package db

import (
	"context"
	"database/sql"
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

// TransferTx performs a money transfer from one account to another
// It creates a new transfer record, updates the balance of two accounts,
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: int64(arg.FromAccountID),
			ToAccountID:   int64(arg.ToAccountID),
			Amount:        int64(arg.Amount),
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: int64(arg.FromAccountID),
			Amount:    -int64(arg.Amount),
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: int64(arg.ToAccountID),
			Amount:    int64(arg.Amount),
		})

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, int64(arg.FromAccountID), -int64(arg.Amount), int64(arg.ToAccountID), int64(arg.Amount))

			if err != nil {
				return err
			}

		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, int64(arg.ToAccountID), int64(arg.Amount), int64(arg.FromAccountID), -int64(arg.Amount))

			if err != nil {
				return err
			}

		}

		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, account1 int64, amount1 int64, account2 int64, amount2 int64) (accountResult1 Account, account2Result Account, err error) {

	accountResult1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1,
		Amount: amount1,
	})

	if err != nil {
		//same thing as return accountResult1{}, accountResult2{}, err
		return
	}

	account2Result, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2,
		Amount: amount2,
	})

	//same thing as return accountResult1, accountResult2, nil
	return

}
