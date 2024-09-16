package main

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/internal/app"
	"github.com/DKhorkov/hmtm-sso/internal/config"
	grpccontroller "github.com/DKhorkov/hmtm-sso/internal/controllers/grpc"
	mocks "github.com/DKhorkov/hmtm-sso/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-sso/internal/services"
	"github.com/DKhorkov/hmtm-sso/internal/usecases"
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
