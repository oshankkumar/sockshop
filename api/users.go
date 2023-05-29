package api

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email"`
	ID        uuid.UUID `json:"id"`
	Links     Links     `json:"_links"`
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

type UserService interface {
	Login(ctx context.Context, username, password string) (*User, error)
	Register(ctx context.Context, user User) (uuid.UUID, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetCard(ctx context.Context, id string) (*Card, error)
	GetUserCards(ctx context.Context, userID string) ([]Card, error)
	GetUserAddresses(ctx context.Context, userID string) ([]Address, error)
	GetAddresses(ctx context.Context, id string) (*Address, error)
	CreateCard(ctx context.Context, card Card, userID string) (uuid.UUID, error)
	CreateAddress(ctx context.Context, addr Address, userID string) (uuid.UUID, error)
}
