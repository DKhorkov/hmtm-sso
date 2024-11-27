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
		settings.Database,
		logger,
	)

	if err != nil {
		panic(err)
	}

	defer dbConnector.CloseConnection()

	usersRepository := repositories.NewCommonUsersRepository(dbConnector)
	authRepository := repositories.NewCommonAuthRepository(dbConnector)
	authService := services.NewCommonAuthService(authRepository, usersRepository)
	usersService := services.NewCommonUsersService(usersRepository)
	useCases := usecases.NewCommonUseCases(
		authService,
		usersService,
		settings.Security.HashCost,
		settings.Security.JWT,
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
