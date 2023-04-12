package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
)

// Store provides all the database operations
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore is a wrapper around queries that provides a set of methods
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {

	// if fn == nil {
	// 	return fmt.Errorf("fn is nil")
	// }

	// if ctx == nil {
	// 	return fmt.Errorf("ctx is nil")
	// }

	// if store.db == nil {
	// 	return fmt.Errorf("db is nil")
	// }

	// if store.Queries == nil {
	// 	return fmt.Errorf("queries is nil")
	// }

	// if store.Queries.db == nil {
	// 	return fmt.Errorf("queries.db is nil")
	// }

	if store == nil || store.db == nil {
		return errors.New("store or store.db is nil")
	}

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
