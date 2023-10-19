package user

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/handlers"
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
		{http.MethodPost, "/login", handlers.LoginHandler(u.userService)},
		{http.MethodPost, "/customers", handlers.RegisterUserHandler(u.userService)},
		{http.MethodGet, "/customers/{id}", handlers.GetUserHandler(u.userService)},
		{http.MethodGet, "/cards/{id}", handlers.GetCardHandler(u.userService)},
		{http.MethodGet, "/addresses/{id}", handlers.GetAddressHandler(u.userService)},
		{http.MethodGet, "/customers/{id}/cards", handlers.GetUserCardsHandler(u.userService)},
		{http.MethodGet, "/customers/{id}/addresses", handlers.GetUserAddressesHandler(u.userService)},
		{http.MethodPost, "/customers/{id}/cards", handlers.CreateCardHandler(u.userService)},
		{http.MethodPost, "/customers/{id}/addresses", handlers.CreateAddressHandler(u.userService)},
	}
	for _, r := range routeDefs {
		mux.Method(r.method, r.pattern, httpkit.ToStdHandler(r.handler))
	}
}
