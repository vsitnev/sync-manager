package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		App   App
		HTTP  HTTP
		Log   Log
		DSN   DSN
	}

	App struct {
		Name    string `env:"APP_NAME"`
		Version string `env:"APP_VERSION"`
		Debug   bool   `env:"APP_DEBUG"`
	}

	HTTP struct {
		Host string `env:"HTTP_HOST"`
		Port string  `env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}

	DSN struct {
		Database string `env:"PG_URL"`
		Amqp     string `env:"AMQP_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse env: %w", err)
	}

	return cfg, nil
}

