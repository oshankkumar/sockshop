package http

import (
	"context"
	"net/http"

	"github.com/oshankkumar/sockshop/api"
)

type sockLister interface {
	ListSocks(ctx context.Context, req *api.ListSockRequest) (*api.ListSockResponse, error)
}

func ListSocksHandler(sockLister sockLister) HTTPHandlerFunc {
	return HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) *Error {
		var req api.ListSockRequest

		resp, err := sockLister.ListSocks(r.Context(), &req)
		if err != nil {
			return &Error{http.StatusInternalServerError, "failed to list socks", err}
		}

		RespondJSON(w, resp, http.StatusOK)
		return nil
	})
}
