package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TrasnferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
}

//store provides all functions to execute db queries and transactions

type SQLStore struct {
	*Queries
	db *sql.DB
}

var txKey = struct{}{}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a database function
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			fmt.Printf("transation error %v, rollback error %v", err, rberr)
		}
		return err
	}
	return tx.Commit()

}

type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TrasnferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID > arg.ToAccountID {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
			if err != nil {
				return err
			}
		}

		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		return nil
	})
	return result, err

}

func addMoney(ctx context.Context, q *Queries, account1ID, account2ID int64, amount1, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2})
	if err != nil {
		return
	}
	return
}
