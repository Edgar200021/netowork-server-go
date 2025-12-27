package testapp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/Edgar200021/netowork-server-go/internal/app"
	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type TestApp struct {
	addressV1 string
	client    *resty.Client
	redis     *redis.Client
	Db        *sql.DB
}

func (a *TestApp) AssertValidationErrors(
	t *testing.T, response *resty.Response,
	keys ...string,
) {
	var respBody struct {
		Errors map[string][]string `json:"errors"`
	}
	err := json.Unmarshal(response.Body(), &respBody)
	require.NoError(t, err)

	for _, key := range keys {
		_, ok := respBody.Errors[key]
		require.True(t, ok)
	}
}

func New(t *testing.T) *TestApp {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Cannot get current file path")
	}

	envPath := path.Join(path.Dir(filename), "../../.env.test")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.New()
	cfg.Application.Port = 0

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Application.Port))
	if err != nil {
		log.Fatal(err)
	}

	client, closeRedis := setupRedis(&cfg.Redis)
	db, closeDb := setupDb(&cfg.Postgres)

	application := app.New(
		ln, cfg, slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{},
			),
		),
	)

	go func() {
		if err := application.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	c := resty.New()

	t.Cleanup(
		func() {
			application.Close(context.Background())
			closeRedis()
			closeDb()
		},
	)

	return &TestApp{
		addressV1: fmt.Sprintf("http://localhost:%d/api/v1", ln.Addr().(*net.TCPAddr).Port),
		client:    c,
		Db:        db,
		redis:     client,
	}
}
