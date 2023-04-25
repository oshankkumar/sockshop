package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/oshankkumar/sockshop/domain"

	"github.com/jmoiron/sqlx"
)

func NewSockStore(db *sqlx.DB) *SockStore {
	return &SockStore{db: db}
}

type SockStore struct {
	db *sqlx.DB
}

func (s *SockStore) List(ctx context.Context, tags []string, order string, limit, offset int) ([]domain.Sock, error) {
	const query = `
	SELECT
		sock.id,
		sock.name,
		sock.description,
		sock.price,
		sock.count,
		sock.image_urls,
		GROUP_CONCAT(tag.name) AS tag_name
	FROM
		sock
	JOIN sock_tag ON
		sock.id = sock_tag.sock_id
	JOIN tag ON
		sock_tag.tag_id = tag.tag_id
	WHERE tag.name IN (?)	
	GROUP BY id ORDER BY ?
	`

	var results []struct {
		ID          string          `db:"id"`
		Name        sql.NullString  `db:"name"`
		Description sql.NullString  `db:"description"`
		Price       sql.NullFloat64 `db:"price"`
		Count       sql.NullInt32   `db:"count"`
		ImageURLs   sql.NullString  `db:"image_urls"`
		TagName     sql.NullString  `db:"tag_name"`
	}

	if err := s.db.SelectContext(ctx, &results, query, tags, order); err != nil {
		return nil, fmt.Errorf("SockStore.List: %w", err)
	}

	var socks []domain.Sock
	for _, res := range results {
		var tags []domain.Tag
		for _, t := range strings.Split(res.TagName.String, ",") {
			tags = append(tags, domain.Tag{Name: t})
		}
		socks = append(socks, domain.Sock{
			ID:          res.ID,
			Name:        res.Name.String,
			Description: res.Description.String,
			ImageURLs:   res.ImageURLs.String,
			Price:       res.Price.Float64,
			Count:       int(res.Count.Int32),
			Tags:        tags,
		})
	}
	return socks, nil
}

func (s *SockStore) Count(ctx context.Context, tags []string) (int, error) {
	panic("implement me")
}

func (s *SockStore) Get(ctx context.Context, id string) (domain.Sock, error) {
	panic("implement me")
}

func (s *SockStore) Tags(ctx context.Context) ([]string, error) {
	panic("implement me")
}

func (s *SockStore) Create(ctx context.Context, sock domain.Sock) error {
	panic("implement me")
}

func (s *SockStore) Update(ctx context.Context, sock domain.Sock) error {
	panic("implement me")
}
