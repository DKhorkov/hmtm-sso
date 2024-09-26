package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/DKhorkov/hmtm-bff/pkg/loadenv"
	"github.com/DKhorkov/hmtm-sso/pkg/logging"
)

func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8070),
		},
		Security: SecurityConfig{
			HashCost: loadenv.GetEnvAsInt("HASH_COST", 8), // Auth speed sensitive if large
			JWT: JWTConfig{
				TTL: time.Hour * time.Duration(
					loadenv.GetEnvAsInt("JWT_TTL", 24),
				),
				Algorithm: loadenv.GetEnv("JWT_ALGORITHM", "HS256"),
				SecretKey: loadenv.GetEnv("JWT_SECRET", "defaultSecret"),
			},
		},
		Databases: DatabasesConfig{
			PostgreSQL: DatabaseConfig{
				Host:         loadenv.GetEnv("POSTGRES_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("POSTGRES_PORT", 5432),
				User:         loadenv.GetEnv("POSTGRES_USER", "postgres"),
				Password:     loadenv.GetEnv("POSTGRES_PASSWORD", "postgres"),
				DatabaseName: loadenv.GetEnv("POSTGRES_DB", "postgres"),
				SSLMode:      loadenv.GetEnv("POSTGRES_SSL_MODE", "disable"),
				Driver:       loadenv.GetEnv("POSTGRES_DRIVER", "postgres"),
			},
		},
		Logging: LoggingConfig{
			Level:       logging.LogLevels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().Format("02-01-2006")),
		},
	}
}

type HTTPConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	SecretKey string
	Algorithm string
	TTL       time.Duration
}

type SecurityConfig struct {
	HashCost int
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
	Driver       string
}

type DatabasesConfig struct {
	PostgreSQL DatabaseConfig
	MySQL      DatabaseConfig
	SQLite     DatabaseConfig
}

type LoggingConfig struct {
	Level       slog.Level
	LogFilePath string
}

type Config struct {
	HTTP      HTTPConfig
	Security  SecurityConfig
	Databases DatabasesConfig
	Logging   LoggingConfig
}
