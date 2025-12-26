package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/Edgar200021/netowork-server-go/internal/app"
	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type TestApp struct {
	address string
}

func New(t *testing.T) *TestApp {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Cannot get current file path")
	}

	envPath := path.Join(path.Dir(filename), "../.env.test")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.New()
	cfg.Application.Port = 0
	cfg.Postgres.Database = fmt.Sprintf("test-%s", uuid.New().String())

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Application.Port))
	if err != nil {
		log.Fatal(err)
	}

	_, closeDb := setupDb(&cfg.Postgres)

	application := app.New(
		ln, cfg, slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{},
			),
		),
	)

	go func() {
		log.Fatal(application.Run())
	}()

	t.Cleanup(
		func() {
			application.Close(context.Background())
			closeDb()
		},
	)

	return &TestApp{
		address: fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port),
	}
}

func setupDb(config *config.PostgresConfig) (*sql.DB, func()) {
	db, err := sql.Open("postgres", config.ConnectionStringWithoutDb())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		fmt.Sprintf(
			"CREATE DATABASE %s",
			pq.QuoteIdentifier(config.Database),
		),
	); err != nil {
		log.Fatal(err)
	}

	connectedDb, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	if err := connectedDb.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Cannot get current file path")
	}

	if err := goose.Up(connectedDb, path.Join(path.Dir(filename), "../migrations")); err != nil {
		log.Fatal(err)
	}

	return connectedDb, func() {
		connectedDb.Close()

		if _, err := db.Exec(
			fmt.Sprintf(
				"DROP DATABASE %s",
				pq.QuoteIdentifier(config.Database),
			),
		); err != nil {
			fmt.Printf("failed to remove database %s\n", err.Error())
		}

		db.Close()
	}
}
