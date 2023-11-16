package api

import "context"

type (
	// Health describes the health of a service
	Health struct {
		Service string `json:"service"`
		Status  string `json:"status"`
		Time    string `json:"time"`
		Details any    `json:"details,omitempty"`
	}

	HealthResponse struct {
		Healths []Health `json:"healths"`
	}
)

type HealthChecker interface {
	CheckHealth(ctx context.Context) ([]Health, error)
}

type HealthCheckerFunc func(ctx context.Context) ([]Health, error)

func (h HealthCheckerFunc) CheckHealth(ctx context.Context) ([]Health, error) { return h(ctx) }
