package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDuplicateUserEntry = errors.New("duplicate entry for user")
)

type DuplicateUserEntryError struct {
	Entity string
	Err    error
}

func (d DuplicateUserEntryError) Error() string { return d.Entity + ":" + d.Err.Error() }

func (d DuplicateUserEntryError) Unwrap() error { return d.Err }

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
	CreateAddress(ctx context.Context, addr *Address, userID string) error
	Delete(ctx context.Context, entity, id string) error
	CreateCard(ctx context.Context, card *Card, id string) error
}

type UserStore interface {
	UserStoreReader
	UserStoreWriter
}
