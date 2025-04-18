package server

import (
	"log/slog"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Add request-related fields to the context
		logAttrs := slog.Group("request_ctx",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("host", r.Host),
		)

		// Log the request details with the context
		slog.Info("Request received", logAttrs)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
