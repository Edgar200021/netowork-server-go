package config

import (
	"fmt"
	"time"
)

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true" validate:"required,hostname|ip"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true" validate:"required,min=1,max=65535"`
	User     string `env:"POSTGRES_USER" env-required:"true" validate:"required,min=1"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true" validate:"required,min=1"`
	Database string `env:"POSTGRES_DB" env-required:"true" validate:"required,min=1"`
	Ssl      bool   `env:"POSTGRES_SSL" env-required:"true" `

	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS" env-default:"4" validate:"gte=1"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS" env-default:"4" validate:"gte=0"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" env-default:"1h" validate:"gt=0"`
	ConnMaxIdleTime time.Duration `env:"POSTGRES_CONN_MAX_IDLE_TIME" env-default:"30m" validate:"gt=0"`
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("%s dbname=%s", c.ConnectionStringWithoutDb(), c.Database)
}

func (c *PostgresConfig) ConnectionStringWithoutDb() string {
	sslMode := "disable"
	if c.Ssl {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		sslMode,
	)
}
