package service

import (
	"context"
	"time"

	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	. "github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	"github.com/Edgar200021/netowork-server-go/pkg/hashing"
	"github.com/google/uuid"
)

func (s *AuthService) ResetPassword(ctx context.Context, data *ResetPasswordRequest) error {

	userId, err := s.cache.GetDel(ctx, s.generateResetPasswordKey(data.Token))
	if err != nil {
		return err
	}

	if userId == "" {
		return ErrResetPasswordTokenInvalidOrExpired
	}

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	user, err := s.repository.GetById(ctx, uuid)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	if user.Email != data.Email {
		return ErrResetPasswordTokenInvalidOrExpired
	}
	if user.IsBanned {
		return ErrUserBanned
	}
	if !user.IsVerified {
		return ErrUserNotVerified
	}

	hashedPassword, err := hashing.HashPassword(data.Password)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	user.Password = hashedPassword
	user.UpdatedAt = &now

	if err := s.repository.Update(ctx, uuid, user); err != nil {
		return err
	}

	return nil
}
