package catalogue

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/handlers"
	"github.com/oshankkumar/sockshop/api/httpkit"
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

func (c *Router) InstallRoutes(mux router.Mux) {
	routeDefs := []struct {
		method  string
		pattern string
		handler httpkit.Handler
	}{
		{http.MethodGet, "/catalogue", handlers.ListSocksHandler(c.catalogueService)},
		{http.MethodGet, "/catalogue/size", handlers.CountTagsHandler(c.sockStore)},
		{http.MethodGet, "/catalogue/{id}", handlers.GetSocksHandler(c.sockStore)},
		{http.MethodGet, "/tags", handlers.TagsHandler(c.sockStore)},
	}
	for _, r := range routeDefs {
		mux.Method(r.method, r.pattern, r.handler)
	}
}

func ImageRouter(path string) router.Router {
	h := http.StripPrefix("/catalogue/images/", http.FileServer(http.Dir(path)))
	hfunc := func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
		h.ServeHTTP(w, r)
		return nil
	}

	return router.RouterFunc(func(mux router.Mux) {
		mux.Method(http.MethodGet, "/catalogue/images/*", httpkit.HandlerFunc(hfunc))
	})
}