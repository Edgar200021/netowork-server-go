package handler

import (
	"net/http"

	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	"github.com/Edgar200021/netowork-server-go/internal/errorhandler"
	"github.com/Edgar200021/netowork-server-go/internal/helpers"

	"github.com/Edgar200021/netowork-server-go/pkg/http_helpers"
)

func (h *AuthHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	var request, err = http_helpers.ParseBody[VerifyAccountRequest](w, r)
	if err != nil {
		return
	}

	res, err := h.authService.VerifyAccount(r.Context(), request)
	if err != nil {
		errorhandler.HandleError(w, err, helpers.LoggerWithRequestId(r.Context(), h.logger))
		return
	}

	h.SetSessionCookie(w, res.SessionId)

	http_helpers.WriteSuccessJson(w, res.UserResponse, http.StatusOK)
}
