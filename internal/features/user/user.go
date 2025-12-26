package user

import (
	"database/sql"
	"log/slog"

	"github.com/Edgar200021/netowork-server-go/internal/features/user/repository"
	"github.com/Edgar200021/netowork-server-go/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

type Feature struct {
	Repository *repository.UserRepository
}

type Dependencies struct {
	Db     *sql.DB
	Logger *slog.Logger
}

func (f *Feature) Handler(middlewares *middlewares.Middlewares) *chi.Mux {

	r := chi.NewMux()
	r.Use(middlewares.Authenticate)

	return r
}

func New(deps *Dependencies) *Feature {

	userRepository := repository.New(deps.Db)

	return &Feature{Repository: userRepository}
}
