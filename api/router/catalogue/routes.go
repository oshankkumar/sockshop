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

func listSocksHandler(sockLister sockLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := sockLister.ListSocks(r.Context(), decodeListReq(r))
		if err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to list socks", Err: err})
			return
		}

		httpkit.RespondJSON(w, resp, http.StatusOK)
	}
}

func countTagsHandler(tagCounter tagCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tags []string
		if tagsval := r.FormValue("tags"); tagsval != "" {
			tags = strings.Split(tagsval, ",")
		}

		c, err := tagCounter.Count(r.Context(), tags)
		if err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to count tags", Err: err})
			return
		}

		httpkit.RespondJSON(w, &api.CountTagsResponse{Size: c}, http.StatusOK)
	}
}

func getSocksHandler(sockGetter sockGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sock, err := sockGetter.Get(r.Context(), chi.URLParam(r, "id"))

		switch {
		case errors.Is(err, domain.ErrNotFound):
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusNotFound, Message: "failed to get sock", Err: err})
			return
		case err != nil:
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to get sock", Err: err})
			return
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
	}
}

func tagsHandler(t tagsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := t.Tags(r.Context())
		if err != nil {
			httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: "failed to get tags", Err: err})
			return
		}

		httpkit.RespondJSON(w, api.TagsResponse{Tags: tags}, http.StatusOK)
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
