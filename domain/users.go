package domain

import (
	"context"

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
	GetAddresses(ctx context.Context) ([]Address, error)
	GetCard(ctx context.Context, id string) (Card, error)
	GetCards(ctx context.Context) ([]Card, error)
}

type UserStoreWriter interface {
	CreateUser(ctx context.Context, user *User) error
	CreateAddress(ctx context.Context, addr *Address, userID string) error
	Delete(ctx context.Context, entity, id string) error
	CreateCard(ctx context.Context, card *Card, id string) error
}

type UserStore interface {
	UserStoreReader
	UserStoreWriter
}
