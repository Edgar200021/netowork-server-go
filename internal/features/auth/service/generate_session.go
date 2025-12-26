package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) GenerateSession(ctx context.Context, userId string) (string, error) {
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	if err := s.cache.Set(
		ctx, s.generateSessionKey(sessionId.String()), userId,
		time.Minute*time.Duration(s.appConfig.SessionTTLMinutes),
	); err != nil {
		return "", err
	}

	return sessionId.String(), nil
}
