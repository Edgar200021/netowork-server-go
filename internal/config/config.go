package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Application ApplicationConfig
	Postgres    PostgresConfig
	Redis       RedisConfig
	Smtp        SmtpConfig
	Logger      LoggerConfig
}

func New() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		log.Fatal(formatValidationError(err))
	}

	return &cfg
}

func formatValidationError(err error) error {
	if ves, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range ves {
			return fmt.Errorf(
				"config validation error: %s.%s (%s)",
				ve.StructNamespace(),
				ve.Field(),
				ve.Tag(),
			)
		}
	}
	return err
}
