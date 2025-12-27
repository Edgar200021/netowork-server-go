package testapp

import (
	"context"
	"strings"

	"github.com/Edgar200021/netowork-server-go/internal/features/auth/constants"
	"github.com/Edgar200021/netowork-server-go/pkg/slice_helpers"
)

func (a *TestApp) GetVerificationToken() (string, error) {
	keys, err := a.redis.Keys(context.Background(), "*").Result()
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	filtered := slice_helpers.Filter(
		keys, func(val string) bool {
			return strings.HasPrefix(val, constants.CacheVerificationPrefix)
		},
	)

	if len(filtered) == 0 {
		return "", nil
	}

	token := strings.Split(filtered[len(filtered)-1], constants.CacheVerificationPrefix)

	return token[1], nil
}
