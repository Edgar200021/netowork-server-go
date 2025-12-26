package helpers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/constants"
	"github.com/Edgar200021/netowork-server-go/internal/dto"
	"github.com/Edgar200021/netowork-server-go/internal/errorhandler"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
)

func UserFromContext(context context.Context, w http.ResponseWriter, logger *slog.Logger) (
	*dto.UserResponse,
	bool,
) {
	userValue := context.Value(constants.UserContext)
	user, ok := userValue.(*dto.UserResponse)
	if !ok {
		errorhandler.HandleError(w, autherrors.ErrUnauthorized, logger)
		return nil, false
	}

	return user, true
}

func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(constants.RequestIDKey).(string); ok {
		return v
	}

	return ""
}

func LoggerWithRequestId(ctx context.Context, logger *slog.Logger) *slog.Logger {
	return logger.With(slog.String("request_id", RequestIDFromContext(ctx)))
}
