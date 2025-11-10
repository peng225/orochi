package psqlutil

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey string

var contextTxKey = txKey("transaction")

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

func (t *Transaction) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	ctx = context.WithValue(ctx, contextTxKey, tx)
	err = fn(ctx)
	if err != nil {
		return fmt.Errorf("function failed: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func TxFromCtx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(contextTxKey).(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}
