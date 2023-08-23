package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Provides all functions to execute queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// Creates a new Store object.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction.
func (store *Store) executeTransaction(context context.Context, function func(*Queries) error) error {
	transaction, error := store.db.BeginTx(context, nil)

	if error != nil {
		return error
	}

	query := New(transaction)
	error = function(query)

	if error != nil {
		rollbackError := transaction.Rollback()

		if rollbackError != nil {
			return fmt.Errorf("Transaction error: %v, rollback error: %v", error, rollbackError)
		}

		return error
	}

	return transaction.Commit()
}

// Contains the input parameters of the transfer transaction.
type TransferTransactionParameters struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Is the result of the transfer transaction.
type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Performs a money transfer from one account to the other.
// Creates a transfer record, add account entries and update accounts balance within a single database transaction.
func (store *Store) TransferTransaction(context context.Context, arguments TransferTransactionParameters) (TransferTransactionResult, error) {
	var result TransferTransactionResult

	error := store.executeTransaction(context, func(queries *Queries) error {
		var error error

		result.Transfer, error = queries.CreateTransfer(context, CreateTransferParams{
			FromAccountID: arguments.FromAccountID,
			ToAccountID:   arguments.ToAccountID,
			Amount:        arguments.Amount,
		})

		if error != nil {
			return error
		}

		result.FromEntry, error = queries.CreateEntry(context, CreateEntryParams{
			AccountID: arguments.FromAccountID,
			Amount:    -arguments.Amount,
		})

		if error != nil {
			return error
		}

		result.ToEntry, error = queries.CreateEntry(context, CreateEntryParams{
			AccountID: arguments.ToAccountID,
			Amount:    arguments.Amount,
		})

		if error != nil {
			return error
		}

		// TODO(augiavedoni): update balances.

		return nil
	})

	return result, error
}
