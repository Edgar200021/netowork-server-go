package handler

import (
	"log/slog"
	"net/http"

	"github.com/Edgar200021/netowork-server-go/internal/config"
	"github.com/Edgar200021/netowork-server-go/internal/features/auth/service"
)

type AuthHandler struct {
	authService *service.AuthService
	appConfig   *config.ApplicationConfig
	logger      *slog.Logger
}

func (h *AuthHandler) SetSessionCookie(w http.ResponseWriter, value string) {
	cookie := &http.Cookie{
		Name:     h.appConfig.SessionCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.appConfig.SessionTTLMinutes * 60,
	}

	http.SetCookie(w, cookie)
}

func (h *AuthHandler) DeleteSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     h.appConfig.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
}

func New(
	authService *service.AuthService, appConfig *config.ApplicationConfig,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService,
		appConfig,
		logger,
	}
}
