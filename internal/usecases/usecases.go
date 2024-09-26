package usecases

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/security"
)

type CommonUseCases struct {
	AuthService  interfaces.AuthService
	UsersService interfaces.UsersService
	HashCost     int
}

func (useCases *CommonUseCases) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	hashedPassword, err := security.HashPassword(userData.Credentials.Password, useCases.HashCost)
	if err != nil {
		return 0, err
	}

	userData.Credentials.Password = hashedPassword
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
