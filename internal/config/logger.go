package config

import "log/slog"

type LoggerConfig struct {
	Level      string `env:"LOG_LEVEL" env-required:"true" validate:"required,oneof=debug info warn error"`
	Structured bool   `env:"LOG_STRUCTURED" env-required:"true"`
	AddSource  bool   `env:"LOG_ADD_SOURCE" env-required:"true"`
}

func (l *LoggerConfig) LoggerOptions() *slog.HandlerOptions {
	var level slog.Level

	switch l.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	return &slog.HandlerOptions{
		Level:     level,
		AddSource: l.AddSource,
	}

}
