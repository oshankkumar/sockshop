package user

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/router"
)

func NewRouter(svc api.UserService) *Router {
	return &Router{userService: svc}
}

type Router struct {
	userService api.UserService
}

func (u *Router) Routes() []router.Route {
	return []router.Route{
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
}
