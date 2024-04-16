package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	sqlx.ExecerContext
	sqlx.PreparerContext
	sqlx.QueryerContext
}

type TxBeginner interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

func RunInTransaction(ctx context.Context, t TxBeginner, runF func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := t.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := runF(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
