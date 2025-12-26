package handler

import (
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/helpers"
	"github.com/Edgar200021/netowork-server-go/pkg/http_helpers"
)

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := helpers.UserFromContext(
		r.Context(), w, helpers.LoggerWithRequestId(r.Context(), h.logger),
	)
	if !ok {
		return
	}

	http_helpers.WriteSuccessJson(w, user, http.StatusOK)
}
