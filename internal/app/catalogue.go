package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/internal/domain"
)

func NewCatalogueService(s domain.SockStore) *CatalogueService {
	return &CatalogueService{sockStore: s}
}

type CatalogueService struct {
	sockStore domain.SockStore
}

func (s *CatalogueService) ListSocks(ctx context.Context, req *api.ListSockParams) (*api.ListSockResponse, error) {
	offset := req.PageSize * (req.PageNum - 1)

	socks, err := s.sockStore.List(ctx, req.Tags, req.Order, req.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("CatalogueService.ListSocks: %w", err)
	}

	var socksResp []api.Sock
	for _, s := range socks {
		var tags []string
		for _, t := range s.Tags {
			tags = append(tags, t.Name)
		}

		socksResp = append(socksResp, api.Sock{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			ImageURL:    strings.Split(s.ImageURLs, ","),
			Price:       s.Price,
			Count:       s.Count,
			Tags:        tags,
		})
	}

	return &api.ListSockResponse{Socks: socksResp}, nil
}
