package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true" validate:"required,hostname|ip"`
	Port     int    `env:"REDIS_PORT" env-required:"true" validate:"required,min=1,max=65535"`
	Password string `env:"REDIS_PASSWORD" env-default:"" validate:"omitempty"`
	Database int    `env:"REDIS_DB" env-default:"0" validate:"gte=0"`
}

func (c *RedisConfig) ConnectOptions() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.Database,
	}
}
