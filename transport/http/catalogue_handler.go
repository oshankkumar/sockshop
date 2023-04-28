package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/domain"
)

type healthChecker interface {
	CheckHealth(ctx context.Context) ([]api.Health, error)
}

type HealthCheckerFunc func(ctx context.Context) ([]api.Health, error)

func (h HealthCheckerFunc) CheckHealth(ctx context.Context) ([]api.Health, error) { return h(ctx) }

type sockLister interface {
	ListSocks(ctx context.Context, req *api.ListSockParams) (*api.ListSockResponse, error)
}

type tagCounter interface {
	Count(ctx context.Context, tags []string) (int, error)
}

type sockGetter interface {
	Get(ctx context.Context, id string) (domain.Sock, error)
}

func HealthCheckHandler(h healthChecker) HTTPHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		hh, err := h.CheckHealth(r.Context())
		if err != nil {
			return &Error{http.StatusInternalServerError, err.Error(), err}
		}

		RespondJSON(w, api.HealthResponse{Healths: hh}, http.StatusOK)
		return nil
	}
}

func ListSocksHandler(sockLister sockLister) HTTPHandlerFunc {
	return HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) *Error {
		resp, err := sockLister.ListSocks(r.Context(), decodeListReq(r))
		if err != nil {
			return &Error{http.StatusInternalServerError, "failed to list socks", err}
		}

		RespondJSON(w, resp, http.StatusOK)
		return nil
	})
}

func CountTagsHandler(tagCounter tagCounter) HTTPHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		var tags []string
		if tagsval := r.FormValue("tags"); tagsval != "" {
			tags = strings.Split(tagsval, ",")
		}

		c, err := tagCounter.Count(r.Context(), tags)
		if err != nil {
			return &Error{http.StatusInternalServerError, "failed to count tags", err}
		}

		RespondJSON(w, &api.CountTagsResponse{Size: c}, http.StatusOK)
		return nil
	}
}

func GetSockHandler(sockGetter sockGetter) HTTPHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		id := mux.Vars(r)["id"]

		sock, err := sockGetter.Get(r.Context(), id)
		if err != nil {
			return &Error{http.StatusInternalServerError, "failed to get sock", err}
		}

		var tags []string
		for _, t := range sock.Tags {
			tags = append(tags, t.Name)
		}

		RespondJSON(w, api.Sock{
			ID:          sock.ID,
			Name:        sock.Name,
			Description: sock.Description,
			ImageURL:    strings.Split(sock.ImageURLs, ","),
			Price:       sock.Price,
			Count:       sock.Count,
			Tags:        tags,
		}, http.StatusOK)

		return nil
	}
}

func decodeListReq(r *http.Request) *api.ListSockParams {
	pageNum := 1
	if page := r.FormValue("page"); page != "" {
		pageNum, _ = strconv.Atoi(page)
	}

	pageSize := 10
	if size := r.FormValue("size"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	order := "id"
	if sort := r.FormValue("sort"); sort != "" {
		order = strings.ToLower(sort)
	}

	var tags []string
	if tagsval := r.FormValue("tags"); tagsval != "" {
		tags = strings.Split(tagsval, ",")
	}

	return &api.ListSockParams{
		Tags:     tags,
		Order:    order,
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}
