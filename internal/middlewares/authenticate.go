package middlewares

import (
	"context"
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/constants"
	"github.com/Edgar200021/netowork-server-go/internal/errorhandler"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	"github.com/Edgar200021/netowork-server-go/internal/helpers"
)

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger := helpers.LoggerWithRequestId(r.Context(), m.logger)

			cookie, err := r.Cookie(m.appConfig.SessionCookieName)
			if err != nil {
				errorhandler.HandleError(
					w, autherrors.ErrUnauthorized, logger,
				)
				return
			}

			user, err := m.authService.Verify(r.Context(), cookie.Value)
			if err != nil {
				errorhandler.HandleError(w, err, logger)
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserContext, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}
