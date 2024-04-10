package api

import (
	"context"

	"github.com/google/uuid"
)

type (
	Sock struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ImageURL    []string  `json:"imageUrl"`
		Price       float64   `json:"price"`
		Count       int       `json:"count"`
		Tags        []string  `json:"tag"`
	}

	ListSockParams struct {
		Tags     []string
		Order    string
		PageNum  int
		PageSize int
	}

	ListSockResponse struct {
		Socks []Sock `json:"sock"`
	}

	CountTagsResponse struct {
		Size int `json:"size"`
	}

	TagsResponse struct {
		Tags []string `json:"tags"`
	}
)

type CatalogueService interface {
	ListSocks(ctx context.Context, req *ListSockParams) (*ListSockResponse, error)
}
