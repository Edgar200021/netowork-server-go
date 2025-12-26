package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/constants"
	"github.com/go-chi/httplog/v3"
	"github.com/google/uuid"
)

func (m *Middlewares) RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := r.Header.Get(constants.RequestIDHeader)

			if requestID == "" {
				if val, err := uuid.NewRandom(); err == nil {
					requestID = val.String()
				}
			}

			if requestID != "" {
				ctx = context.WithValue(r.Context(), constants.RequestIDKey, requestID)
				r = r.WithContext(ctx)
				httplog.SetAttrs(ctx, slog.String("request_id", requestID))
			}

			next.ServeHTTP(w, r)
		},
	)
}
