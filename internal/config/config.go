package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/stdlib-ticket-system/pkg/utils"
	"github.com/joho/godotenv"
)

const (
	ContextTimeout       = time.Second * 10
	AccessTokenDuration  = time.Hour * 1
	RefreshTokenDuration = time.Hour * 24 * 7
)

type EnvConfig struct {
	DB  DBConfig  `envPrefix:"DB_"`
	JWT JWTConfig `envPrefix:"JWT_"`
}

type DBConfig struct {
	DBUser    string `env:"USER" validate:"required"`
	DBPass    string `env:"PASSWORD" validate:"required"`
	DBName    string `env:"NAME" validate:"required"`
	DBHost    string `env:"HOST" envDefault:"localhost"`
	DBPort    int    `env:"PORT" envDefault:"5432"`
	DBSSLMode string `env:"SSL_MODE" envDefault:"disable"`
}

type JWTConfig struct {
	SecretKey  string `env:"SECRET_KEY" validate:"required"`
	RefreshKey string `env:"REFRESH_KEY" validate:"required"`
}

func LoadConfig(path string) (*EnvConfig, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("load env failed: %w", err)
	}

	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}

	if err := utils.Validate(cfg); err != nil {
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
