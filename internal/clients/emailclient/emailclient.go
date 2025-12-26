package emailclient

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/Edgar200021/netowork-server-go/internal/config"
)

type EmailClient struct {
	auth      smtp.Auth
	addr      string
	from      string
	appConfig *config.ApplicationConfig
}

func New(cfg *config.SmtpConfig, appCfg *config.ApplicationConfig) (*EmailClient, error) {
	auth, err := сheckSMTP(cfg)
	if err != nil {
		return nil, err
	}

	return &EmailClient{
		auth:      auth,
		addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		from:      cfg.User,
		appConfig: appCfg,
	}, nil
}

func сheckSMTP(cfg *config.SmtpConfig) (smtp.Auth, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	c, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("smtp dial failed: %w", err)
	}
	defer c.Close()

	if err := c.Hello(cfg.Host); err != nil {
		return nil, fmt.Errorf("smtp hello failed: %w", err)
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(
			&tls.Config{
				ServerName: cfg.Host,
			},
		); err != nil {
			return nil, fmt.Errorf("starttls failed: %w", err)
		}

	}

	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)

	if err := c.Auth(auth); err != nil {
		return nil, fmt.Errorf("smtp auth failed: %w", err)
	}
	return auth, nil
}
