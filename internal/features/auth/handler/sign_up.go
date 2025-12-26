package handler

import (
	"net/http"

	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	"github.com/Edgar200021/netowork-server-go/internal/errorhandler"
	"github.com/Edgar200021/netowork-server-go/internal/helpers"
	"github.com/Edgar200021/netowork-server-go/pkg/http_helpers"
)

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var request, err = http_helpers.ParseBody[SignUpRequest](w, r)
	if err != nil {
		return
	}

	if err := h.authService.SignUp(r.Context(), request); err != nil {
		errorhandler.HandleError(w, err, helpers.LoggerWithRequestId(r.Context(), h.logger))
		return
	}

	http_helpers.WriteSuccessJson(w, "Success", http.StatusCreated)
}
