package main

import (
	"github.com/DKhorkov/hmtm-sso/internal/app"
	grpccontroller "github.com/DKhorkov/hmtm-sso/internal/controllers/grpc"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	mocks "github.com/DKhorkov/hmtm-sso/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-sso/internal/services"
	"github.com/DKhorkov/hmtm-sso/internal/usecases"

	"github.com/DKhorkov/hmtm-sso/internal/config"
)

func main() {
	// logger := logging.GetInstance(logging.LogLevels.DEBUG)
	settings := config.New()

	usersRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*entities.User{}}
	authRepository := usersRepository
	authService := &services.CommonAuthService{
		AuthRepository:  authRepository,
		UsersRepository: usersRepository,
	}

	usersService := &services.CommonUsersService{
		UsersRepository: usersRepository,
	}

	useCases := &usecases.CommonUseCases{
		AuthService:  authService,
		UsersService: usersService,
	}

	controller := grpccontroller.New(settings.HTTP.Host, settings.HTTP.Port, useCases)

	application := app.New(controller)
	application.Run()
}
