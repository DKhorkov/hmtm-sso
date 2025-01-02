package main

import (
	"github.com/DKhorkov/hmtm-sso/internal/app"
	"github.com/DKhorkov/hmtm-sso/internal/config"
	grpccontroller "github.com/DKhorkov/hmtm-sso/internal/controllers/grpc"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
	"github.com/DKhorkov/hmtm-sso/internal/services"
	"github.com/DKhorkov/hmtm-sso/internal/usecases"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
)

func main() {
	settings := config.New()
	logger := logging.GetInstance(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	dbConnector, err := db.New(
		db.BuildDsn(settings.Database),
		settings.Database.Driver,
		logger,
		// db.WithMaxIdleConnections(1),
		// db.WithMaxConnectionLifetime(time.Second*2),
		// db.WithMaxConnectionIdleTime(time.Second*5),
		//db.WithMaxOpenConnections(1),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = dbConnector.Close(); err != nil {
			logging.LogError(logger, "Failed to close db connections pool", err)
		}
	}()

	usersRepository := repositories.NewCommonUsersRepository(dbConnector, logger)
	usersService := services.NewCommonUsersService(
		usersRepository,
		logger,
	)

	authRepository := repositories.NewCommonAuthRepository(dbConnector, logger)
	authService := services.NewCommonAuthService(
		authRepository,
		usersRepository,
		logger,
	)

	useCases := usecases.NewCommonUseCases(
		authService,
		usersService,
		settings.Security,
	)

	controller := grpccontroller.New(
		settings.HTTP.Host,
		settings.HTTP.Port,
		useCases,
		logger,
	)

	application := app.New(controller)
	application.Run()
}
