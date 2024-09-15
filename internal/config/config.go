package config

import "github.com/DKhorkov/hmtm-bff/pkg/loadenv"

func New() *Config {
	return &Config{
		HTTP: HTTPConfigs{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8070),
		},
	}
}

type HTTPConfigs struct {
	Host string
	Port int
}

type Config struct {
	HTTP HTTPConfigs
}
