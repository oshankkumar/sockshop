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

func loginHandler(l loginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, pass, ok := r.BasicAuth()
		if !ok {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusUnauthorized, Message: "user not authorised", Err: api.ErrUnauthorized})
			return
		}

		user, err := l.Login(r.Context(), username, pass)

		switch {
		case errors.Is(err, api.ErrUnauthorized):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusUnauthorized, Message: "user not authorised", Err: err})
			return
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user not found", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "user login failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, user, http.StatusOK)
	}
}

func registerUserHandler(ur userRegisterationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user api.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err})
			return
		}

		id, err := ur.Register(r.Context(), user)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusConflict, Message: "username or email already exists", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "user registeration failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, api.CreateResponse{ID: id}, http.StatusCreated)
	}
}

func getUserHandler(us userGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user not exist", Err: api.ErrNotFound})
			return
		}

		user, err := us.GetUser(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user not found", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "get user details failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, user, http.StatusOK)
	}
}

func getCardHandler(cg cardGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardID := chi.URLParam(r, "id")
		if cardID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "card does not exist", Err: api.ErrNotFound})
		}

		card, err := cg.GetCard(r.Context(), cardID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "card not found", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "get card details failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, card, http.StatusOK)
	}
}

func getUserCardsHandler(cg userCardsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user does not exist", Err: api.ErrNotFound})
		}

		cards, err := cg.GetUserCards(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "card not found", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "get card details failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, api.UserCardsResponse{Cards: cards}, http.StatusOK)
	}
}

func getAddressHandler(ag addressGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addrID := chi.URLParam(r, "id")
		if addrID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "address does not exist", Err: api.ErrNotFound})
			return
		}

		addr, err := ag.GetAddresses(r.Context(), addrID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "address not found", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "get address failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, addr, http.StatusOK)
	}
}

func getUserAddressesHandler(ag userAddressesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user does not exist", Err: api.ErrNotFound})
			return
		}

		addrs, err := ag.GetUserAddresses(r.Context(), userID)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "address not found", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "get address failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, api.UserAdressesResponse{Addresses: addrs}, http.StatusOK)
	}
}

func createCardHandler(cc cardCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user not exists", Err: api.ErrNotFound})
			return
		}

		var card api.Card
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err})
			return
		}

		id, err := cc.CreateCard(r.Context(), card, userID)

		switch {
		case errors.As(err, &domain.DuplicateEntryError{}):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusConflict, Message: "card already registered", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "adding card failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, api.CreateResponse{ID: id}, http.StatusOK)
	}
}

func createAddressHandler(ac addressCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "user not exists", Err: api.ErrNotFound})
			return
		}

		var addr api.Address
		if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusBadRequest, Message: "json unmarshal failed", Err: err})
			return
		}

		id, err := ac.CreateAddress(r.Context(), addr, userID)
		if err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "card creation failed", Err: err})
			return
		}

		httpkit.RespondJSON(w, api.CreateResponse{ID: id}, http.StatusOK)
	}
}
