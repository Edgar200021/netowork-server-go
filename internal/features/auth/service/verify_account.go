package service

import (
	"context"
	"time"

	. "github.com/Edgar200021/netowork-server-go/internal/dto"
	. "github.com/Edgar200021/netowork-server-go/internal/features/auth/autherrors"
	"github.com/google/uuid"
)

func (s *AuthService) VerifyAccount(
	ctx context.Context, data *VerifyAccountRequest,
) (*VerifyAccountResponse, error) {
	userId, err := s.cache.GetDel(ctx, s.generateVerificationKey(data.Token))
	if err != nil {
		return nil, err
	}

	if userId == "" {
		return nil, ErrVerificationTokenInvalidOrExpired
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
		return nil, ErrUserNotFound
	}

	if user.Email != data.Email {
		return nil, ErrVerificationTokenInvalidOrExpired
	}
	if user.IsBanned {
		return nil, ErrUserBanned
	}
	if user.IsVerified {
		return nil, ErrUserAlreadyVerified
	}

	now := time.Now().UTC()

	user.IsVerified = true
	user.UpdatedAt = &now

	if err := s.repository.Update(
		ctx, uuid, user,
	); err != nil {
		return nil, err
	}

	sessionId, err := s.GenerateSession(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &VerifyAccountResponse{
		UserResponse: UserResponse{
			ID:         user.ID.String(),
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
			Email:      user.Email,
			Role:       user.Role,
			Balance:    user.Balance,
			IsVerified: true,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		SessionId: sessionId,
	}, nil
}
