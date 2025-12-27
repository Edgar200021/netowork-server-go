package testapp

import (
	"context"
	"fmt"
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

func (a *TestApp) GetResetPasswordToken() (string, error) {
	keys, err := a.redis.Keys(context.Background(), "*").Result()
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	filtered := slice_helpers.Filter(
		keys, func(val string) bool {
			return strings.HasPrefix(val, constants.CacheResetPasswordPrefix)
		},
	)

	if len(filtered) == 0 {
		return "", nil
	}

	token := strings.Split(filtered[len(filtered)-1], constants.CacheResetPasswordPrefix)

	return token[1], nil
}

func (a *TestApp) ExpireVerificationToken() error {
	token, err := a.GetVerificationToken()
	if err != nil {
		return err
	}

	if _, err := a.redis.Expire(
		context.Background(), fmt.Sprintf(
			"%s%s",
			constants.CacheVerificationPrefix, token,
		), 0,
	).Result(); err != nil {
		return err
	}

	return nil
}
