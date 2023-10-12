package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/kellemNegasi/bank-system/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	st := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := fmt.Sprint(util.RandInt(500, 800))
	errs := make(chan error)
	results := make(chan TransferTxResult)
	tests := 5

	// run the transactions in a separate go routines.
	for i := 0; i < tests; i++ {
		go func() {
			txResult, err := st.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- txResult
		}()
	}

	for i := 0; i < tests; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// transfer check
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = st.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, "-"+amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = st.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = st.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO: check account balance

	}

}
