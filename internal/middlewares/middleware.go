package middlewares

import (
	"log/slog"

	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/service"
)

type Middlewares struct {
	appConfig   *config.ApplicationConfig
	authService *service.AuthService
	logger      *slog.Logger
}

func New(
	appConfig *config.ApplicationConfig,
	authService *service.AuthService,
	logger *slog.Logger,
) *Middlewares {

	return &Middlewares{
		appConfig,
		authService,
		logger,
	}
}
