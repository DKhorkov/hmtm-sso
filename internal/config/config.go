package config

import (
	"time"

	"github.com/DKhorkov/hmtm-bff/pkg/loadenv"
)

func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8070),
		},
		Security: SecurityConfig{
			HashCost: loadenv.GetEnvAsInt("HASH_COST", 14),
			JWT: JWTConfig{
				TTL: time.Hour * time.Duration(
					loadenv.GetEnvAsInt("JWT_TTL", 24),
				),
				Algorithm: loadenv.GetEnv("JWT_ALGORITHM", "SHA256"),
				SecretKey: loadenv.GetEnv("JWT_SECRET", "defaultSecret"),
			},
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

type Config struct {
	HTTP     HTTPConfig
	Security SecurityConfig
}
