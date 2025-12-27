package service

import (
	"context"
	"time"

	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	"github.com/google/uuid"
)

func (s *AuthService) Authenticate(ctx context.Context, sessionId string) (*UserResponse, error) {
	userId, err := s.cache.GetEx(
		ctx, s.generateSessionKey(sessionId),
		time.Minute*time.Duration(s.appConfig.SessionTTLMinutes),
	)

	if err != nil {
		return nil, err
	}

	if userId == "" {
		return nil, autherrors.ErrUnauthorized
	}

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.GetById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, autherrors.ErrUnauthorized
	}

	if !user.IsVerified {
		return nil, autherrors.ErrUserNotVerified
	}

	if user.IsBanned {
		return nil, autherrors.ErrUserBanned
	}

	return &UserResponse{
		ID:         user.ID.String(),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Email:      user.Email,
		Role:       user.Role,
		Balance:    user.Balance,
		IsVerified: user.IsVerified,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	}, nil
}
