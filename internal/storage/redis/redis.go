package redis

import (
	"context"

	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(config *config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(config.ConnectOptions())

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Redis{
		client,
	}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
