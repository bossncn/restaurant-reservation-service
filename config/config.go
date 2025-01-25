package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv string `envconfig:"APP_ENV" validate:"required" default:"development"`
}

func (c *Config) Validate() error {
	validate := validator.New()
	// Perform validation
	err := validate.Struct(c)

	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	var config Config

	_ = godotenv.Load()
	_ = envconfig.Process("", &config)

	// Validate the configuration
	if err := config.Validate(); err != nil {
		fmt.Println("Configuration validation failed:", err)
		return nil, err
	}

	return &config, nil
}
