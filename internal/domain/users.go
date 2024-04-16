package domain

import (
	"context"

	"github.com/oshankkumar/sockshop/internal/db"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	Email      string    `db:"email"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	Salt       string    `db:"salt"`
	AddressIDs []string  `db:"-"`
	CardIDs    []string  `db:"-"`
}

type Address struct {
	ID       uuid.UUID `db:"id"`
	Street   string    `db:"street"`
	Number   string    `db:"number"`
	Country  string    `db:"country"`
	City     string    `db:"city"`
	PostCode string    `db:"postcode"`
}

type Card struct {
	ID      uuid.UUID `db:"id"`
	LongNum string    `db:"long_num"`
	Expires string    `db:"expires"`
	CCV     string    `db:"ccv"`
}

type UserStoreReader interface {
	GetUserByName(ctx context.Context, uname string) (User, error)
	GetUser(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetAddress(ctx context.Context, id string) (Address, error)
	GetUserAddresses(ctx context.Context, userID string) ([]Address, error)
	GetCard(ctx context.Context, id string) (Card, error)
	GetUserCards(ctx context.Context, userID string) ([]Card, error)
}

type UserStoreWriter interface {
	CreateUser(ctx context.Context, user *User) error
	CreateAddress(ctx context.Context, addrID string, userID string) error
	Delete(ctx context.Context, entity, id string) error
	CreateCard(ctx context.Context, cardID string, id string) error
}

type CardStore interface {
	CreateCard(ctx context.Context, card *Card) error
	WithTx(db db.DB) CardStore
}

type AddressStore interface {
	CreateAddress(ctx context.Context, addr *Address) error
	WithTx(db db.DB) AddressStore
}

type UserStore interface {
	WithTx(db db.DB) UserStore
	UserStoreReader
	UserStoreWriter
}
