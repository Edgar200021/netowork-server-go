package app

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Edgar200021/netowork-server-go/internal/clients/emailclient"
	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth"
	"github.com/Edgar200021/netowork-server-go/internal/features/user"
	"github.com/Edgar200021/netowork-server-go/internal/middlewares"
	"github.com/Edgar200021/netowork-server-go/internal/storage/postgres"
	"github.com/Edgar200021/netowork-server-go/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

type App struct {
	Run   func() error
	Close func(ctx context.Context)
}

func New(listener net.Listener, config *config.Config, logger *slog.Logger) *App {
	db, err := postgres.New(config.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	redis, err := redis.New(&config.Redis)
	if err != nil {
		log.Fatal(err)
	}

	emailSender, err := emailclient.New(&config.Smtp, &config.Application)
	if err != nil {
		log.Fatal(err)
	}

	userFeature := user.New(
		&user.Dependencies{
			Db:     db,
			Logger: logger,
		},
	)
	authFeature := auth.New(
		&auth.Dependencies{
			UserRepository: userFeature.Repository,
			EmailSender:    emailSender,
			Cache:          redis,
			AppConfig:      &config.Application,
			Logger:         logger,
		},
	)

	middlewares := middlewares.New(&config.Application, authFeature.Service, logger)

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(httplog.RequestLogger(logger, nil))
	r.Use(middlewares.RequestId)
	r.Use(middleware.Recoverer)

	r.Route(
		"/api/v1", func(r chi.Router) {
			r.Mount("/auth", authFeature.Handler(middlewares))
			r.Mount("/user", userFeature.Handler(middlewares))
		},
	)

	server := http.Server{
		Handler:      r,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return &App{
		Run: func() error {
			return server.Serve(listener)
		},
		Close: func(ctx context.Context) {
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Error shutting down server", "error", err)
			}

			if err := redis.Close(); err != nil {
				logger.Error("Error closing Redis", "error", err)
			}

			if err := db.Close(); err != nil {
				logger.Error("Error closing Postgres", "error", err)
			}
		},
	}
}
