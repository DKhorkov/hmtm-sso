package config

import "github.com/DKhorkov/hmtm-bff/pkg/loadenv"

func GetConfig() *Config {
	return &Config{
		GRPC: GRPCConfigs{
			Port: loadenv.GetEnvAsInt("GRPC_PORT", 8070),
		},
	}
}

type GRPCConfigs struct {
	Port int
}

type Config struct {
	GRPC GRPCConfigs
}
