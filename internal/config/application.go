package config

type ApplicationConfig struct {
	Port                          int    `env:"APP_PORT" env-required:"true" validate:"required,min=1,max=65535"`
	AccountVerificationTTLMinutes int    `env:"APP_ACCOUNT_VERIFICATION_TTL_MINUTES" env-required:"true" validate:"required,gt=0,lte=1440"`
	ResetPasswordTTLMinutes       int    `env:"APP_RESET_PASSWORD_TTL_MINUTES" env-required:"true" validate:"required,gt=5,lte=10"`
	SessionTTLMinutes             int    `env:"APP_SESSION_TTL_MINUTES" env-required:"true" validate:"required,gt=1400,lte=43200"`
	SessionCookieName             string `env:"APP_SESSION_COOKIE_NAME" env-required:"true" validate:"required,min=1"`
	ClientUrl                     string `env:"APP_CLIENT_URL" env-required:"true" validate:"required,url"`
	AccountVerificationPath       string `env:"APP_CLIENT_ACCOUNT_VERIFICATION_PATH" env-required:"true" validate:"required,min=1"`
	ResetPasswordPath             string `env:"APP_CLIENT_RESET_PASSWORD_PATH" env-required:"true" validate:"required,min=1"`
}
