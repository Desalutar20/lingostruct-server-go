package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Database    DatabseConfig
	Application ApplicationConfig
	Redis       RedisConfig
	Smtp        SmtpConfig
	RateLimit   RateLimitConfig
}

func New() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		ve := err.(validator.ValidationErrors)

		for _, e := range ve {
			fmt.Printf(
				"config error: %s\n",
				e.Error(),
			)
		}

		panic("config validation failed")
	}

	return &cfg
}
