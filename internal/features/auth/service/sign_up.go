package service

import (
	"context"
	"time"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	"github.com/Edgar200021/netowork-server-go/pkg/hashing"
	token2 "github.com/Edgar200021/netowork-server-go/pkg/token"
)

func (s *AuthService) SignUp(
	ctx context.Context,
	data *SignUpRequest,
) error {
	user, err := s.repository.GetByEmail(
		ctx,
		data.Email,
	)
	if err != nil {
		return err
	}

	if user != nil {
		return autherrors.ErrUserAlreadyExists
	}

	hashedPassword, err := hashing.HashPassword(data.Password)
	if err != nil {
		return err
	}

	id, err := s.repository.Create(
		ctx, model.Users{
			Email:     data.Email,
			Password:  hashedPassword,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Role:      data.Role,
		},
	)

	if err != nil {
		return err
	}

	token := token2.GenerateSecureToken(20)

	if err := s.cache.Set(
		ctx, s.generateVerificationKey(token), id,
		time.Minute*time.Duration(s.appConfig.AccountVerificationTTLMinutes),
	); err != nil {
		return err
	}

	if err := s.emailSender.SendVerificationEmail(token, data.Email); err != nil {
		return err
	}

	return nil
}
