package user

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/httpkit"
	"github.com/oshankkumar/sockshop/api/router"
)

func NewRouter(svc api.UserService) *Router {
	return &Router{userService: svc}
}

type Router struct {
	userService api.UserService
}

func (u *Router) InstallRoutes(mux router.Mux) {
	routeDefs := []struct {
		method  string
		pattern string
		handler httpkit.Handler
	}{
		{http.MethodPost, "/login", loginHandler(u.userService)},
		{http.MethodPost, "/customers", registerUserHandler(u.userService)},
		{http.MethodGet, "/customers/{id}", getUserHandler(u.userService)},
		{http.MethodGet, "/cards/{id}", getCardHandler(u.userService)},
		{http.MethodGet, "/addresses/{id}", getAddressHandler(u.userService)},
		{http.MethodGet, "/customers/{id}/cards", getUserCardsHandler(u.userService)},
		{http.MethodGet, "/customers/{id}/addresses", getUserAddressesHandler(u.userService)},
		{http.MethodPost, "/customers/{id}/cards", createCardHandler(u.userService)},
		{http.MethodPost, "/customers/{id}/addresses", createAddressHandler(u.userService)},
	}
	for _, r := range routeDefs {
		mux.Method(r.method, r.pattern, r.handler)
	}
}
