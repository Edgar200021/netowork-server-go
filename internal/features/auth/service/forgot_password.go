package service

import (
	"context"
	"time"

	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	token2 "github.com/Edgar200021/netowork-server-go/pkg/token"
)

func (s *AuthService) ForgotPassword(ctx context.Context, data *ForgotPasswordRequest) error {
	user, err := s.repository.GetByEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	if user == nil {
		return autherrors.ErrUserNotFound
	}
	if !user.IsVerified {
		return autherrors.ErrUserNotVerified
	}
	if user.IsBanned {
		return autherrors.ErrUserBanned
	}

	token := token2.GenerateSecureToken(16)
	if err := s.cache.Set(
		ctx, s.generateResetPasswordKey(token), user.ID.String(),
		time.Minute*time.Duration(s.appConfig.ResetPasswordTTLMinutes),
	); err != nil {
		return err
	}

	if err := s.emailSender.SendResetPasswordEmail(token, user.Email); err != nil {
		return err
	}

	return nil

}
