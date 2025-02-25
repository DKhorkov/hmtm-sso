package main

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	customnats "github.com/DKhorkov/libs/nats"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-sso/internal/app"
	"github.com/DKhorkov/hmtm-sso/internal/config"
	grpccontroller "github.com/DKhorkov/hmtm-sso/internal/controllers/grpc"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
	"github.com/DKhorkov/hmtm-sso/internal/services"
	"github.com/DKhorkov/hmtm-sso/internal/usecases"
)

func main() {
	settings := config.New()
	logger := logging.New(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	dbConnector, err := db.New(
		db.BuildDsn(settings.Database),
		settings.Database.Driver,
		logger,
		db.WithMaxOpenConnections(settings.Database.Pool.MaxOpenConnections),
		db.WithMaxIdleConnections(settings.Database.Pool.MaxIdleConnections),
		db.WithMaxConnectionLifetime(settings.Database.Pool.MaxConnectionLifetime),
		db.WithMaxConnectionIdleTime(settings.Database.Pool.MaxConnectionIdleTime),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = dbConnector.Close(); err != nil {
			logging.LogError(logger, "Failed to close db connections pool", err)
		}
	}()

	traceProvider, err := tracing.New(settings.Tracing.Server)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = traceProvider.Shutdown(context.Background()); err != nil {
			logging.LogError(logger, "Error shutting down tracer", err)
		}
	}()

	natsPublisher, err := customnats.NewPublisher(
		settings.NATS.ClientURL,
		nats.Name(settings.NATS.Publisher.Name),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = natsPublisher.Close(); err != nil {
			logging.LogError(logger, "Failed to close nats publisher", err)
		}
	}()

	usersRepository := repositories.NewUsersRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Users,
	)

	usersService := services.NewUsersService(
		usersRepository,
		logger,
	)

	authRepository := repositories.NewAuthRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Auth,
	)

	authService := services.NewAuthService(
		authRepository,
		usersRepository,
		logger,
	)

	useCases := usecases.New(
		authService,
		usersService,
		settings.Security,
		settings.Validation,
		natsPublisher,
		settings.NATS,
		logger,
	)

	controller := grpccontroller.New(
		settings.HTTP.Host,
		settings.HTTP.Port,
		useCases,
		logger,
		traceProvider,
		settings.Tracing.Spans.Root,
	)

	application := app.New(controller)
	application.Run()
}
