package handlers

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/api/httpkit"
)

func HealthCheckHandler(h api.HealthChecker) httpkit.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
		hh, err := h.CheckHealth(r.Context())
		if err != nil {
			return &httpkit.Error{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
		}

		httpkit.RespondJSON(w, api.HealthResponse{Healths: hh}, http.StatusOK)
		return nil
	}
}
