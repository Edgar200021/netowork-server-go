package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *Redis) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	value, err := r.client.GetEx(ctx, key, expiration).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", err
	}

	return value, nil
}
