package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/kellemNegasi/bank-system/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	st := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account1Balance, err := decimal.NewFromString(account1.Balance)
	require.NoError(t, err)
	account2Balance, err := decimal.NewFromString(account2.Balance)
	require.NoError(t, err)
	amount := fmt.Sprint(util.RandInt(10, 20))
	fmt.Printf(">> Before balance1 = %s, balance2= %s, amount= %v", account1.Balance, account2.Balance, amount)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	amountDecimal, err := decimal.NewFromString(amount)
	require.NoError(t, err)

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

	existed := map[int64]bool{}
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

		// check fromAccount

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Get the decimal from the string format
		fromBalance, err := decimal.NewFromString(fromAccount.Balance)
		require.NoError(t, err)
		toBalance, err := decimal.NewFromString(toAccount.Balance)
		require.NoError(t, err)
		diff1 := account1Balance.Sub(fromBalance)
		diff2 := toBalance.Sub(account2Balance)

		rem := diff1.BigInt().Int64() % amountDecimal.BigInt().Int64()
		require.Equal(t, diff1, diff2)
		require.True(t, diff1.GreaterThan(decimal.Zero))
		require.True(t, rem == 0)

		k := diff1.Div(amountDecimal).BigInt().Int64()
		require.True(t, k >= 1 && k <= int64(tests))
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check for updated accounts

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedBalance1, err := decimal.NewFromString(updatedAccount1.Balance)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	updatedBalance2, err := decimal.NewFromString(updatedAccount2.Balance)
	require.NoError(t, err)

	a1 := amountDecimal.Mul(decimal.NewFromInt(int64(tests)))
	require.Equal(t, account1Balance.Sub(a1), updatedBalance1)
	require.Equal(t, account2Balance.Add(a1), updatedBalance2)

	fmt.Printf(">> After balance1 = %s, balance2= %s", updatedAccount1.Balance, updatedAccount2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {
	st := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account1Balance, err := decimal.NewFromString(account1.Balance)
	require.NoError(t, err)
	account2Balance, err := decimal.NewFromString(account2.Balance)
	require.NoError(t, err)
	amount := fmt.Sprint(util.RandInt(10, 20))
	fmt.Printf(">> Before balance1 = %s, balance2= %s, amount= %v", account1.Balance, account2.Balance, amount)
	errs := make(chan error)

	tests := 10

	// run the transactions in a separate go routines.
	for i := 0; i < tests; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := st.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < tests; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check for updated accounts

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedBalance1, err := decimal.NewFromString(updatedAccount1.Balance)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	updatedBalance2, err := decimal.NewFromString(updatedAccount2.Balance)
	require.NoError(t, err)

	require.Equal(t, account1Balance, updatedBalance1)
	require.Equal(t, account2Balance, updatedBalance2)

	fmt.Printf(">> After balance1 = %s, balance2= %s", updatedAccount1.Balance, updatedAccount2.Balance)

}
