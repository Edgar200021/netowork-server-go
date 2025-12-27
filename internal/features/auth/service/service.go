package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	"github.com/Edgar200021/netowork-server-go/internal/config"
	. "github.com/Edgar200021/netowork-server-go/internal/features/auth/constants"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user model.Users) (string, error)
	GetByEmail(ctx context.Context, email string) (*model.Users, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Users, error)
	Update(ctx context.Context, id uuid.UUID, user *model.Users) error
}

type EmailSender interface {
	SendVerificationEmail(token, to string) error
	SendResetPasswordEmail(token, to string) error
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetEx(ctx context.Context, key string, expiration time.Duration) (string, error)
	GetDel(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type AuthService struct {
	repository  UserRepository
	emailSender EmailSender
	cache       Cache
	appConfig   *config.ApplicationConfig
}

func (s *AuthService) generateVerificationKey(token string) string {
	return fmt.Sprintf("%s%s", CacheVerificationPrefix, token)
}

func (s *AuthService) generateSessionKey(sessionId string) string {
	return fmt.Sprintf("%s%s", CacheSessionPrefix, sessionId)
}

func (s *AuthService) generateResetPasswordKey(token string) string {
	return fmt.Sprintf("%s%s", CacheResetPasswordPrefix, token)
}

func New(
	repository UserRepository, emailSender EmailSender, cache Cache,
	appConfig *config.ApplicationConfig,
) *AuthService {
	return &AuthService{
		repository,
		emailSender,
		cache,
		appConfig,
	}
}
