package api

import "context"

type (
	Sock struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ImageURL    []string `json:"imageUrl"`
		Price       float32  `json:"price"`
		Count       int      `json:"count"`
		Tags        []string `json:"tag"`
	}

	ListSockRequest struct {
		Tags     []string
		Order    string
		PageNum  int
		PageSize int
	}

	ListSockResponse struct {
		Socks []Sock `json:"sock"`
	}

	CatalogueService interface {
		ListSocks(ctx context.Context, req *ListSockRequest) (*ListSockResponse, error)
	}
)
