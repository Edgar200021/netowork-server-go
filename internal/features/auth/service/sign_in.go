package service

import (
	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	. "github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"

	"github.com/Edgar200021/netowork-server-go/pkg/hashing"

	"context"
)

func (s *AuthService) SignIn(
	ctx context.Context, data *SignInRequest,
) (*SignInResponse, error) {
	user, err := s.repository.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !hashing.VerifyPassword(data.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	if user.IsBanned {
		return nil, ErrUserBanned
	}
	if !user.IsVerified {
		return nil, ErrUserNotVerified
	}

	sessionId, err := s.GenerateSession(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return &SignInResponse{
		UserResponse: UserResponse{
			ID:         user.ID.String(),
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
			Email:      user.Email,
			Role:       user.Role,
			Balance:    user.Balance,
			IsVerified: user.IsVerified,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		SessionId: sessionId,
	}, nil
}
