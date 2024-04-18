package catalogue

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/router"
	"github.com/oshankkumar/sockshop/internal/domain"
)

func NewRouter(cs api.CatalogueService, ss domain.SockStore) *Router {
	return &Router{catalogueService: cs, sockStore: ss}
}

type Router struct {
	catalogueService api.CatalogueService
	sockStore        domain.SockStore
}

func (c *Router) Routes() []router.Route {
	return []router.Route{
		{Method: http.MethodGet, Pattern: "/catalogue", Handler: listSocksHandler(c.catalogueService)},
		{Method: http.MethodGet, Pattern: "/catalogue/size", Handler: countTagsHandler(c.sockStore)},
		{Method: http.MethodGet, Pattern: "/catalogue/{id}", Handler: getSocksHandler(c.sockStore)},
		{Method: http.MethodGet, Pattern: "/tags", Handler: tagsHandler(c.sockStore)},
	}
}

func ImageRouter(path string) router.RouterFunc {
	return func() []router.Route {
		return []router.Route{
			{
				Method:  http.MethodGet,
				Pattern: "/catalogue/images/*",
				Handler: http.StripPrefix("/catalogue/images/", http.FileServer(http.Dir(path))),
			},
		}
	}
}
