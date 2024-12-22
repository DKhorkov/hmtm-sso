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
	defer func() {
		if err = usersRepository.Close(); err != nil {
			logging.LogError(logger, "Failed to close Users repository", err)
		}
	}()

	usersService := services.NewCommonUsersService(
		usersRepository,
		logger,
	)

	authRepository := repositories.NewCommonAuthRepository(dbConnector)
	defer func() {
		if err = authRepository.Close(); err != nil {
			logging.LogError(logger, "Failed to close Auth repository", err)
		}
	}()

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
