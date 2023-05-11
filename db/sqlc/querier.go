// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteTransfer(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	//method normal
	// UPDATE users SET
	//   username = $2,
	//   hashed_password = $3,
	//   full_name = $4,
	//   email = $5
	// WHERE username = $1
	// RETURNING *;
	// method when not all fields are updated
	// UPDATE users
	//   Set
	//   hashed_password = CASE
	//     WHEN @set_hashed_password::boolean = TRUE THEN @hashed_password
	//     ELSE hashed_password
	//   END,
	//   full_name = CASE
	//     WHEN @set_full_name = TRUE THEN @full_name
	//     ELSE full_name
	//   END,
	//   email = CASE
	//     WHEN  @set_email = TRUE THEN @email
	//     ELSE email
	//   END
	// WHERE username = @username
	// RETURNING *;
	// method when not all fields are updated method 3
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
