package app

import grpcapp "github.com/DKhorkov/hmtm-sso/internal/app/grpc"

type App struct {
	GRPC *grpcapp.GrpcApp
}

func New(grpcHost string, grpcPort int) *App {
	grpcApp := grpcapp.New(grpcHost, grpcPort)
	return &App{
		GRPC: grpcApp,
	}
}
