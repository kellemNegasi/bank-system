package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

// Store extends the functionalities of Queries and adds the capability of running SQL transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore returns a Store object.
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// executeTx executes the provided callback function inside a database transaction.
func (st *Store) executeTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := st.db.Begin()
	if err != nil {
		return err
	}

	queries := New(tx)

	err = fn(queries)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return fmt.Errorf("rollback error : %v, transaction error: %v", err, rollBackErr)
		}

		return err
	}

	return tx.Commit()
}

// TransferTxParams defines the collection of parameters needed to execute a transfer transaction.
type TransferTxParams struct {
	FromAccountID int64  `from_account_id`
	ToAccountID   int64  `to_account_id`
	Amount        string `amount`
}

// TransferTxResult holds the result of transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `transfer`
	FromAccount Account  `from_account`
	ToAccount   Account  `to_account`
	FromEntry   Entry    `from_entry`
	ToEntry     Entry    `to_entry`
}

// TransferTx executes the transfer transaction.
func (st *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var txResult TransferTxResult

	cbFunc := func(q *Queries) error {
		var err error
		TxParams := CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		}
		txResult.Transfer, err = q.CreateTransfer(ctx, TxParams)
		if err != nil {
			return err
		}

		fromEntryParams := CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    "-" + args.Amount,
		}

		toEntryParams := CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		}

		txResult.FromEntry, err = q.CreateEntry(ctx, fromEntryParams)
		if err != nil {
			return err
		}

		txResult.ToEntry, err = q.CreateEntry(ctx, toEntryParams)
		if err != nil {
			return err
		}

		// Update accounts balance.

		amountDec, err := decimal.NewFromString(args.Amount)
		if err != nil {
			return err
		}

		// Here let's first update the account with the smaller ID.
		// This is for predictable order to avoid deadlocks.
		if args.FromAccountID < args.ToAccountID {

			// Update account1 First
			txResult.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
				ID:     args.FromAccountID,
				Amount: "-" + amountDec.String(),
			})

			if err != nil {
				return err
			}

			// Update account2
			txResult.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
				ID:     args.ToAccountID,
				Amount: amountDec.String(),
			})
			if err != nil {
				return err
			}
		} else {

			// Update account2 First
			txResult.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
				ID:     args.ToAccountID,
				Amount: amountDec.String(),
			})
			if err != nil {
				return err
			}

			// Update account1
			txResult.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
				ID:     args.FromAccountID,
				Amount: "-" + amountDec.String(),
			})

			if err != nil {
				return err
			}

		}

		return nil
	}

	err := st.executeTx(ctx, cbFunc)
	return txResult, err
}
