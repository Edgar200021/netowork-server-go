package config

type SmtpConfig struct {
	Host     string `env:"SMTP_HOST" env-required:"true" validate:"required,hostname|ip"`
	Port     int    `env:"SMTP_PORT" env-required:"true" validate:"required,min=1,max=65535"`
	User     string `env:"SMTP_USER" env-required:"true" validate:"required,min=1"`
	Password string `env:"SMTP_PASSWORD" env-required:"true" validate:"required,min=1"`
}
