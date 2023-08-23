package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, error := store.TransferTransaction(context.Background(), TransferTransactionParameters{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errors <- error
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		error := <-errors
		require.NoError(t, error)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer.

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, error = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, error)

		// Check entries.

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, error = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, error)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, error = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, error)

		// TODO(augiavedoni): check accounts balance.
	}
}
