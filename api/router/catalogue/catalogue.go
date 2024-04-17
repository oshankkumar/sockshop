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
		{http.MethodGet, "/catalogue", listSocksHandler(c.catalogueService)},
		{http.MethodGet, "/catalogue/size", countTagsHandler(c.sockStore)},
		{http.MethodGet, "/catalogue/{id}", getSocksHandler(c.sockStore)},
		{http.MethodGet, "/tags", tagsHandler(c.sockStore)},
	}
}

func ImageRouter(path string) router.RouterFunc {
	return func() []router.Route {
		return []router.Route{
			{http.MethodGet, "/catalogue/images/*", http.StripPrefix("/catalogue/images/", http.FileServer(http.Dir(path)))},
		}
	}
}
