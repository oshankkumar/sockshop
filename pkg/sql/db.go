package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	sqlx.ExecerContext
	sqlx.PreparerContext
	sqlx.QueryerContext
}

type TxBeginer interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type DBTx interface {
	TxBeginer
	DB
}

func RunInTx(ctx context.Context, opts *sql.TxOptions, db TxBeginer, run func(context.Context, DB) error) error {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("db.RunInTx Begin Tx: %w", err)
	}

	defer tx.Rollback()

	if err := run(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db.RunInTx Commit Tx: %w", err)
	}

	return nil
}
