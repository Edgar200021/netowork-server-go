package handler

import (
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/errorhandler"
	"github.com/Edgar200021/netowork-server-go/internal/helpers"
	"github.com/Edgar200021/netowork-server-go/pkg/http_helpers"
)

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := helpers.LoggerWithRequestId(r.Context(), h.logger)
	cookie, err := r.Cookie(h.appConfig.SessionCookieName)
	if err != nil {
		errorhandler.HandleError(w, err, logger)
		return
	}

	if err := h.authService.Logout(r.Context(), cookie.Value); err != nil {
		errorhandler.HandleError(w, err, logger)
		return
	}

	h.DeleteSessionCookie(w)
	http_helpers.WriteSuccessJson(w, "Success", http.StatusOK)
}
