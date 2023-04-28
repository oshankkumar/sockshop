package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/oshankkumar/sockshop/domain"
	sqli "github.com/oshankkumar/sockshop/pkg/sql"
)

const baseQuery = "SELECT sock.id, sock.name, sock.description, sock.price, sock.count, sock.image_urls, " +
	"GROUP_CONCAT(tag.name) AS tag_name " +
	"FROM sock JOIN sock_tag ON sock.id = sock_tag.sock_id JOIN tag ON sock_tag.tag_id = tag.id "

func NewSockStore(db sqli.DB) *SockStore {
	return &SockStore{db: db}
}

type SockStore struct {
	db sqli.DB
}

func (s *SockStore) WithTx(db sqli.DB) domain.SockStore {
	return &SockStore{db: db}
}

func (s *SockStore) List(ctx context.Context, tags []string, order string, limit, offset int) ([]domain.Sock, error) {
	var tagCond string
	if len(tags) > 0 {
		tagCond = "WHERE tag.name IN ( ?" + strings.Repeat(",?", len(tags)-1) + ") "
	}

	query := baseQuery + tagCond + "GROUP BY id ORDER BY ?"

	var results []struct {
		ID          string          `db:"id"`
		Name        sql.NullString  `db:"name"`
		Description sql.NullString  `db:"description"`
		Price       sql.NullFloat64 `db:"price"`
		Count       sql.NullInt32   `db:"count"`
		ImageURLs   sql.NullString  `db:"image_urls"`
		TagName     sql.NullString  `db:"tag_name"`
	}

	args := make([]interface{}, 0, len(tags)+1)
	for _, t := range tags {
		args = append(args, t)
	}
	args = append(args, order)

	if err := sqlx.SelectContext(ctx, s.db, &results, query, args...); err != nil {
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
	query := "SELECT COUNT(DISTINCT sock.sock_id) FROM sock JOIN sock_tag ON sock.sock_id=sock_tag.sock_id JOIN tag ON sock_tag.tag_id=tag.tag_id"

	var tagCond string
	if len(tags) > 0 {
		tagCond = "WHERE tag.name IN ( ?" + strings.Repeat(",?", len(tags)-1) + ") "
	}

	args := make([]interface{}, 0, len(tags))
	for _, t := range tags {
		args = append(args, t)
	}

	query += tagCond + ";"

	var count int
	if err := sqlx.GetContext(ctx, s.db, &count, query, args...); err != nil {
		return 0, fmt.Errorf("SockStore.Count: %w", err)
	}

	return count, nil
}

func (s *SockStore) Get(ctx context.Context, id string) (domain.Sock, error) {
	query := baseQuery + " WHERE sock.sock_id =? GROUP BY sock.sock_id;"

	var result struct {
		ID          string          `db:"id"`
		Name        sql.NullString  `db:"name"`
		Description sql.NullString  `db:"description"`
		Price       sql.NullFloat64 `db:"price"`
		Count       sql.NullInt32   `db:"count"`
		ImageURLs   sql.NullString  `db:"image_urls"`
		TagName     sql.NullString  `db:"tag_name"`
	}

	if err := sqlx.GetContext(ctx, s.db, &result, query, id); err != nil {
		return domain.Sock{}, fmt.Errorf("SockStore.Get(%s): %w", id, err)
	}

	var tags []domain.Tag
	for _, t := range strings.Split(result.TagName.String, ",") {
		tags = append(tags, domain.Tag{Name: t})
	}

	return domain.Sock{
		ID:          result.ID,
		Name:        result.Name.String,
		Description: result.Description.String,
		ImageURLs:   result.ImageURLs.String,
		Price:       result.Price.Float64,
		Count:       int(result.Count.Int32),
		Tags:        tags,
	}, nil
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
