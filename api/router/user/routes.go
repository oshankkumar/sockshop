package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/httpkit"
	"github.com/oshankkumar/sockshop/internal/domain"
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

func loginHandler(l loginService) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		username, pass, ok := r.BasicAuth()
		if !ok {
			return &httpkit.Error{Code: http.StatusUnauthorized, Message: "user not authorised", Err: api.ErrUnauthorized}

		}

		user, err := l.Login(r.Context(), username, pass)

		switch {
		case errors.Is(err, api.ErrUnauthorized):
			return &httpkit.Error{Code: http.StatusUnauthorized, Message: "user not authorised", Err: err}
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user not found", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "user login failed", Err: err}
		}

		httpkit.RespondJSON(w, user, http.StatusOK)
		return nil
	}
}

func registerUserHandler(ur userRegisterationService) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var user api.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return &httpkit.Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err}
		}

		id, err := ur.Register(r.Context(), user)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			return &httpkit.Error{Code: http.StatusConflict, Message: "username or email already exists", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "user registeration failed", Err: err}
		}

		httpkit.RespondJSON(w, api.CreateResponse{ID: id}, http.StatusCreated)
		return nil
	}
}

func getUserHandler(us userGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user not exist", Err: api.ErrNotFound}
		}

		user, err := us.GetUser(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user not found", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "get user details failed", Err: err}
		}

		httpkit.RespondJSON(w, user, http.StatusOK)
		return nil
	}
}

func getCardHandler(cg cardGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		cardID := chi.URLParam(r, "id")
		if cardID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "card does not exist", Err: api.ErrNotFound}
		}

		card, err := cg.GetCard(r.Context(), cardID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "card not found", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "get card details failed", Err: err}
		}

		httpkit.RespondJSON(w, card, http.StatusOK)
		return nil
	}
}

func getUserCardsHandler(cg userCardsGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user does not exist", Err: api.ErrNotFound}
		}

		cards, err := cg.GetUserCards(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "card not found", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "get card details failed", Err: err}
		}

		httpkit.RespondJSON(w, api.UserCardsResponse{Cards: cards}, http.StatusOK)
		return nil
	}
}

func getAddressHandler(ag addressGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		addrID := chi.URLParam(r, "id")
		if addrID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "address does not exist", Err: api.ErrNotFound}
		}

		addr, err := ag.GetAddresses(r.Context(), addrID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "address not found", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "get address failed", Err: err}
		}

		httpkit.RespondJSON(w, addr, http.StatusOK)
		return nil
	}
}

func getUserAddressesHandler(ag userAddressesGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user does not exist", Err: api.ErrNotFound}
		}

		addrs, err := ag.GetUserAddresses(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "address not found", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "get address failed", Err: err}
		}

		httpkit.RespondJSON(w, api.UserAdressesResponse{Addresses: addrs}, http.StatusOK)
		return nil
	}
}

func createCardHandler(cc cardCreator) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user not exists", Err: api.ErrNotFound}
		}

		var card api.Card
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			return &httpkit.Error{Code: http.StatusBadRequest, Message: "invalid json body in request", Err: err}
		}

		id, err := cc.CreateCard(r.Context(), card, userID)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			return &httpkit.Error{Code: http.StatusConflict, Message: "card already registered", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "adding card failed", Err: err}
		}

		httpkit.RespondJSON(w, api.CreateResponse{ID: id}, http.StatusOK)
		return nil
	}
}

func createAddressHandler(ac addressCreator) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			return &httpkit.Error{Code: http.StatusNotFound, Message: "user not exists", Err: api.ErrNotFound}
		}

		var addr api.Address
		if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
			return &httpkit.Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err}
		}

		id, err := ac.CreateAddress(r.Context(), addr, userID)
		if err != nil {
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "address creation failed", Err: err}
		}

		httpkit.RespondJSON(w, api.CreateResponse{ID: id}, http.StatusOK)
		return nil
	}
}
