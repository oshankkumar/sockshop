package catalogue

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/httpkit"
	"github.com/oshankkumar/sockshop/internal/domain"
)

type sockLister interface {
	ListSocks(ctx context.Context, req *api.ListSockParams) (*api.ListSockResponse, error)
}

type tagCounter interface {
	Count(ctx context.Context, tags []string) (int, error)
}

type sockGetter interface {
	Get(ctx context.Context, id string) (domain.Sock, error)
}

type tagsGetter interface {
	Tags(ctx context.Context) ([]string, error)
}

func listSocksHandler(sockLister sockLister) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
		resp, err := sockLister.ListSocks(r.Context(), decodeListReq(r))
		if err != nil {
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to list socks", Err: err}
		}

		httpkit.RespondJSON(w, resp, http.StatusOK)
		return nil
	}
}

func countTagsHandler(tagCounter tagCounter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
		var tags []string
		if tagsval := r.FormValue("tags"); tagsval != "" {
			tags = strings.Split(tagsval, ",")
		}

		c, err := tagCounter.Count(r.Context(), tags)
		if err != nil {
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to count tags", Err: err}
		}

		httpkit.RespondJSON(w, &api.CountTagsResponse{Size: c}, http.StatusOK)
		return nil
	}
}

func getSocksHandler(sockGetter sockGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
		sock, err := sockGetter.Get(r.Context(), chi.URLParam(r, "id"))

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &httpkit.Error{Code: http.StatusNotFound, Message: "failed to get sock", Err: err}
		case err != nil:
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to get sock", Err: err}
		}

		var tags []string
		for _, t := range sock.Tags {
			tags = append(tags, t.Name)
		}

		httpkit.RespondJSON(w, api.Sock{
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

func tagsHandler(t tagsGetter) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
		tags, err := t.Tags(r.Context())
		if err != nil {
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to get tags", Err: err}
		}

		httpkit.RespondJSON(w, api.TagsResponse{Tags: tags}, http.StatusOK)
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
