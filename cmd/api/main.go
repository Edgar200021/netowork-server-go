package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Edgar200021/netowork-server-go/internal/app"
	"github.com/Edgar200021/netowork-server-go/internal/config"
	logger "github.com/Edgar200021/netowork-server-go/internal/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cfg := config.New()
	appLogger := logger.New(&cfg.Logger)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Application.Port))
	if err != nil {
		log.Fatal(err)
	}

	application := app.New(ln, cfg, appLogger)

	appLogger.Info(fmt.Sprintf("Server running on port %d", cfg.Application.Port))
	go func() {
		if err := application.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			application.Close(context.Background())
			log.Fatal(err)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	application.Close(ctx)
}
