package auth

import (
	"database/sql"
	"log/slog"

	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/handler"
	service "github.com/Edgar200021/netowork-server-go/internal/features/auth/service"
	"github.com/Edgar200021/netowork-server-go/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

type Feature struct {
	Service     *service.AuthService
	handler     *handler.AuthHandler
	middlewares *middlewares.Middlewares
}

type Dependencies struct {
	UserRepository service.UserRepository
	EmailSender    service.EmailSender
	Cache          service.Cache
	AppConfig      *config.ApplicationConfig
	Db             *sql.DB
	Logger         *slog.Logger
}

func (f *Feature) Handler(middlewares *middlewares.Middlewares) *chi.Mux {

	r := chi.NewMux()

	r.Post("/sign-up", f.handler.SignUp)
	r.Post("/sign-in", f.handler.SignIn)
	r.Post("/verify-account", f.handler.VerifyAccount)
	r.Post("/forgot-password", f.handler.ForgotPassword)
	r.Post("/reset-password", f.handler.ResetPassword)

	r.Group(
		func(r chi.Router) {
			r.Use(middlewares.Authenticate)
			r.Post("/logout", f.handler.Logout)
			r.Get("/me", f.handler.GetMe)
		},
	)

	return r
}

func New(deps *Dependencies) *Feature {
	authService := service.New(
		deps.UserRepository, deps.EmailSender, deps.Cache, deps.AppConfig,
	)
	authHandler := handler.New(authService, deps.AppConfig, deps.Logger)

	return &Feature{
		handler: authHandler,
		Service: authService,
	}
}
