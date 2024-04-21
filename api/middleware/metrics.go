package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/oshankkumar/sockshop/api/httpkit"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Byte = 1 << (10 * iota)
	KB
	MB
)

func WithMetrics() httpkit.MiddlewareFunc {
	var (
		reqCount = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "sockshop",
			Name:      "http_request",
			Help:      "The total number of http request served",
		}, []string{"method", "pattern", "code"})

		reqLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "sockshop",
			Name:      "http_request_latency_ms",
			Help:      "http request latency in millisecond",
			Buckets:   []float64{10, 20, 50, 100, 500, 1000, 2000},
		}, []string{"method", "pattern", "code"})

		respSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "sockshop",
			Name:      "http_resp_size_bytes",
			Help:      "http response size in bytes",
			Buckets:   []float64{32 * Byte, 64 * Byte, 128 * Byte, 256 * Byte, 512 * Byte, 1 * KB, 2 * KB, 4 * KB, 8 * KB},
		}, []string{"method", "pattern"})
	)

	return func(method, pattern string, h httpkit.Handler) httpkit.Handler {
		return httpkit.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()

			wr, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				wr = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			err := h.ServeHTTP(wr, r)

			took := time.Since(start).Milliseconds()

			code := strconv.Itoa(wr.Status())

			reqCount.WithLabelValues(method, pattern, code).Inc()
			reqLatency.WithLabelValues(method, pattern, code).Observe(float64(took))
			respSize.WithLabelValues(method, pattern).Observe(float64(wr.BytesWritten()))
			return err
		})
	}
}
