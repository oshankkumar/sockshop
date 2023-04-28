package domain

import (
	"context"

	"github.com/oshankkumar/sockshop/pkg/sql"
)

type Tag struct {
	ID   string
	Name string
}

type Sock struct {
	ID          string
	Name        string
	Description string
	ImageURLs   string
	Price       float64
	Count       int
	Tags        []Tag
}

type SockStore interface {
	SockStoreTransctioner
	SockStoreReader
	SockStoreWriter
}

type SockStoreTransctioner interface {
	WithTx(sql.DB) SockStore
}

type SockStoreReader interface {
	List(ctx context.Context, tags []string, order string, limit, offset int) ([]Sock, error)
	Count(ctx context.Context, tags []string) (int, error)
	Get(ctx context.Context, id string) (Sock, error)
	Tags(ctx context.Context) ([]string, error)
}

type SockStoreWriter interface {
	Create(ctx context.Context, sock Sock) error
	Update(ctx context.Context, sock Sock) error
}
