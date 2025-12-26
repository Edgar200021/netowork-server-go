package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func (r *Redis) Get(ctx context.Context, key string) (string, error) {

	value, err := r.client.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", err
	}

	return value, nil
}
