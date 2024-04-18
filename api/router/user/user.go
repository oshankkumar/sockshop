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
		{Method: http.MethodPost, Pattern: "/login", Handler: loginHandler(u.userService)},
		{Method: http.MethodPost, Pattern: "/customers", Handler: registerUserHandler(u.userService)},
		{Method: http.MethodGet, Pattern: "/customers/{id}", Handler: getUserHandler(u.userService)},
		{Method: http.MethodGet, Pattern: "/cards/{id}", Handler: getCardHandler(u.userService)},
		{Method: http.MethodGet, Pattern: "/addresses/{id}", Handler: getAddressHandler(u.userService)},
		{Method: http.MethodGet, Pattern: "/customers/{id}/cards", Handler: getUserCardsHandler(u.userService)},
		{Method: http.MethodGet, Pattern: "/customers/{id}/addresses", Handler: getUserAddressesHandler(u.userService)},
		{Method: http.MethodPost, Pattern: "/customers/{id}/cards", Handler: createCardHandler(u.userService)},
		{Method: http.MethodPost, Pattern: "/customers/{id}/addresses", Handler: createAddressHandler(u.userService)},
	}
}
