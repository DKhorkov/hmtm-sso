package config

import "github.com/DKhorkov/hmtm-bff/pkg/loadenv"

func GetConfig() *Config {
	return &Config{
		GRPC: GRPCConfigs{
			Host: loadenv.GetEnv("GRPC_HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("GRPC_PORT", 8070),
		},
	}
}

type GRPCConfigs struct {
	Host string
	Port int
}

type Config struct {
	GRPC GRPCConfigs
}
