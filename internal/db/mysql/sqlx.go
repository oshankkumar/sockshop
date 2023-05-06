package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/oshankkumar/sockshop/domain"
)

func GetContext(ctx context.Context, q sqlx.QueryerContext, dest interface{}, query string, args ...interface{}) error {
	err := sqlx.GetContext(ctx, q, dest, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func SelectContext(ctx context.Context, q sqlx.QueryerContext, dest interface{}, query string, args ...interface{}) error {
	err := sqlx.SelectContext(ctx, q, dest, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}
