package usecases

import (
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

type CommonUseCases struct {
	AuthService  interfaces.AuthService
	UsersService interfaces.UsersService
}

func (useCases *CommonUseCases) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return useCases.AuthService.RegisterUser(userData)
}

func (useCases *CommonUseCases) LoginUser(userData entities.LoginUserDTO) (string, error) {
	return useCases.AuthService.LoginUser(userData)
}

func (useCases *CommonUseCases) GetUserByID(id int) (*entities.User, error) {
	return useCases.UsersService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]*entities.User, error) {
	return useCases.UsersService.GetAllUsers()
}
