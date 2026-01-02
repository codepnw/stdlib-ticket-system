package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/stdlib-ticket-system/pkg/utils"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DB DBConfig `envPrefix:"DB_"`
}

type DBConfig struct {
	DBUser    string `env:"USER" validate:"required"`
	DBPass    string `env:"PASSWORD" validate:"required"`
	DBName    string `env:"NAME" validate:"required"`
	DBHost    string `env:"HOST" envDefault:"localhost"`
	DBPort    int    `env:"PORT" envDefault:"5432"`
	DBSSLMode string `env:"SSL_MODE" envDefault:"disable"`
}

func LoadConfig(path string) (*EnvConfig, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("load env failed: %w", err)
	}
	
	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}
	
	if err := utils.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate env failed: %w", err)
	}
	return cfg, nil
}

func (cfg *EnvConfig) GetDBConnection() string {
	db := cfg.DB
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.DBUser,
		db.DBPass,
		db.DBHost,
		db.DBPort,
		db.DBName,
		db.DBSSLMode,
	)
}