package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/domain"
)

type loginService interface {
	Login(ctx context.Context, username, password string) (*api.User, error)
}

type userRegisterationService interface {
	Register(ctx context.Context, user api.User) (uuid.UUID, error)
}

type userGetter interface {
	GetUser(ctx context.Context, id string) (*api.User, error)
}

type cardGetter interface {
	GetCard(ctx context.Context, id string) (*api.Card, error)
}

type userCardsGetter interface {
	GetUserCards(ctx context.Context, userID string) ([]api.Card, error)
}

type userAddressesGetter interface {
	GetUserAddresses(ctx context.Context, userID string) ([]api.Address, error)
}

type addressGetter interface {
	GetAddresses(ctx context.Context, id string) (*api.Address, error)
}

type cardCreator interface {
	CreateCard(ctx context.Context, card api.Card, userID string) (uuid.UUID, error)
}

type addressCreator interface {
	CreateAddress(ctx context.Context, addr api.Address, userID string) (uuid.UUID, error)
}

func LoginHandler(l loginService) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		username, pass, ok := r.BasicAuth()
		if !ok {
			return &Error{Code: http.StatusUnauthorized, Message: "user not authorised", Err: api.ErrUnauthorized}
		}

		user, err := l.Login(r.Context(), username, pass)

		switch {
		case errors.Is(err, api.ErrUnauthorized):
			return &Error{http.StatusUnauthorized, "user not authorised", err}
		case errors.Is(err, domain.ErrNotFound):
			return &Error{http.StatusNotFound, "user not found", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "user login failed", err}
		}

		RespondJSON(w, user, http.StatusOK)
		return nil
	}
}

func RegisterUserHandler(ur userRegisterationService) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		var user api.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return &Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err}
		}

		id, err := ur.Register(r.Context(), user)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			return &Error{http.StatusConflict, "username or email already exists", err}
		case err != nil:
			return &Error{Code: http.StatusInternalServerError, Message: "user registeration failed", Err: err}
		}

		RespondJSON(w, api.CreateResponse{ID: id}, http.StatusCreated)
		return nil
	}
}

func GetUserHandler(us userGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		userID := mux.Vars(r)["id"]
		if userID == "" {
			return &Error{Code: http.StatusNotFound, Message: "user not exist", Err: api.ErrNotFound}
		}

		user, err := us.GetUser(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &Error{http.StatusNotFound, "user not found", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "get user details failed", err}
		}

		RespondJSON(w, user, http.StatusOK)
		return nil
	}
}

func GetCardHandler(cg cardGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		cardID := mux.Vars(r)["id"]
		if cardID == "" {
			return &Error{Code: http.StatusNotFound, Message: "card does not exist", Err: api.ErrNotFound}
		}

		card, err := cg.GetCard(r.Context(), cardID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &Error{http.StatusNotFound, "card not found", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "get card details failed", err}
		}

		RespondJSON(w, card, http.StatusOK)
		return nil
	}
}

func GetUserCardsHandler(cg userCardsGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		userID := mux.Vars(r)["id"]
		if userID == "" {
			return &Error{Code: http.StatusNotFound, Message: "user does not exist", Err: api.ErrNotFound}
		}

		cards, err := cg.GetUserCards(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &Error{http.StatusNotFound, "card not found", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "get card details failed", err}
		}

		RespondJSON(w, api.UserCardsResponse{Cards: cards}, http.StatusOK)
		return nil
	}
}

func GetAddressHandler(ag addressGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		addrID := mux.Vars(r)["id"]
		if addrID == "" {
			return &Error{Code: http.StatusNotFound, Message: "address does not exist", Err: api.ErrNotFound}
		}

		addr, err := ag.GetAddresses(r.Context(), addrID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &Error{http.StatusNotFound, "address not found", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "get address failed", err}
		}

		RespondJSON(w, addr, http.StatusOK)
		return nil
	}
}

func GetUserAddressesHandler(ag userAddressesGetter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		userID := mux.Vars(r)["id"]
		if userID == "" {
			return &Error{Code: http.StatusNotFound, Message: "user does not exist", Err: api.ErrNotFound}
		}

		addrs, err := ag.GetUserAddresses(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &Error{http.StatusNotFound, "address not found", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "get address failed", err}
		}

		RespondJSON(w, api.UserAdressesResponse{Addresses: addrs}, http.StatusOK)
		return nil
	}
}

func CreateCardHandler(cc cardCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		userID := mux.Vars(r)["id"]
		if userID == "" {
			return &Error{Code: http.StatusNotFound, Message: "user not exists", Err: api.ErrNotFound}
		}

		var card api.Card
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			return &Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err}
		}

		id, err := cc.CreateCard(r.Context(), card, userID)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			return &Error{http.StatusConflict, "card already registered", err}
		case err != nil:
			return &Error{http.StatusInternalServerError, "adding card failed", err}
		}

		RespondJSON(w, api.CreateResponse{ID: id}, http.StatusOK)
		return nil
	}
}

func CreateAddressHandler(ac addressCreator) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		userID := mux.Vars(r)["id"]
		if userID == "" {
			return &Error{Code: http.StatusNotFound, Message: "user not exists", Err: api.ErrNotFound}
		}

		var addr api.Address
		if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
			return &Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err}
		}

		id, err := ac.CreateAddress(r.Context(), addr, userID)
		if err != nil {
			return &Error{http.StatusInternalServerError, "card creation failed", err}
		}

		RespondJSON(w, api.CreateResponse{ID: id}, http.StatusOK)
		return nil
	}
}
