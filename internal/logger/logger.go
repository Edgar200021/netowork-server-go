package logger

import (
	"log/slog"
	"os"

	"github.com/Edgar200021/netowork-server-go/internal/config"
)

func New(config *config.LoggerConfig) *slog.Logger {
	var handler slog.Handler

	if config.Structured {
		handler = slog.NewJSONHandler(os.Stdout, config.LoggerOptions())
	} else {
		handler = slog.NewTextHandler(os.Stdout, config.LoggerOptions())
	}

	return slog.New(handler)
}
