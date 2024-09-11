package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DKhorkov/hmtm-sso/internal/app"

	"github.com/DKhorkov/hmtm-sso/internal/config"
)

func main() {
	// logger := logging.GetInstance(logging.LogLevels.DEBUG)
	settings := config.GetConfig()
	application := app.New(settings.GRPC.Host, settings.GRPC.Port)

	// Launch asynchronous for graceful shutdown purpose:
	go application.GRPC.Run()

	// Graceful shutdown. When system signal will be received, signal.Notify function will write it to channel.
	// After this event, main goroutine will be unblocked (<-stopChannel blocks it) and application will be
	// gracefully stopped:
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	<-stopChannel
	application.GRPC.Stop()
}
