package errorhandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	"github.com/Edgar200021/netowork-server-go/pkg/http_helpers"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

var errorCodeMap = map[error]int{
	autherrors.ErrUserAlreadyExists:                  http.StatusBadRequest,
	autherrors.ErrVerificationTokenInvalidOrExpired:  http.StatusBadRequest,
	autherrors.ErrUserAlreadyVerified:                http.StatusBadRequest,
	autherrors.ErrUserBanned:                         http.StatusBadRequest,
	autherrors.ErrUserNotVerified:                    http.StatusBadRequest,
	autherrors.ErrInvalidCredentials:                 http.StatusBadRequest,
	autherrors.ErrResetPasswordTokenInvalidOrExpired: http.StatusBadRequest,
	autherrors.ErrUserNotFound:                       http.StatusNotFound,
	autherrors.ErrUnauthorized:                       http.StatusUnauthorized,
	ErrUnauthorized:                                  http.StatusUnauthorized,
}

func HandleError(w http.ResponseWriter, err error, log *slog.Logger) {
	for e, code := range errorCodeMap {
		if errors.Is(err, e) {
			http_helpers.WriteErrorJson(w, e.Error(), code)
			return
		}
	}

	log.Error("internal error", slog.Any("error", err))

	http_helpers.WriteErrorJson(w, "Something went wrong", http.StatusInternalServerError)
}
