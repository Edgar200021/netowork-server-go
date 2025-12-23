package config

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

func New() *Config {
	return &Config{}
}
