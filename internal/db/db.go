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

type TxBeginer interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type DBTx interface {
	TxBeginer
	DB
}
