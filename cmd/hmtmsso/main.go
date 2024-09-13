package main

import (
	"github.com/DKhorkov/hmtm-sso/internal/app"
	grpccontroller "github.com/DKhorkov/hmtm-sso/internal/controllers/grpc"

	"github.com/DKhorkov/hmtm-sso/internal/config"
)

func main() {
	// logger := logging.GetInstance(logging.LogLevels.DEBUG)
	settings := config.GetConfig()
	controller := grpccontroller.New(settings.GRPC.Host, settings.GRPC.Port)
	application := app.New(controller)
	application.Run()
}
