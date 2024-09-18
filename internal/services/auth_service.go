package services

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/security"
)

type CommonAuthService struct {
	AuthRepository  interfaces.AuthRepository
	UsersRepository interfaces.UsersRepository
}

func (service *CommonAuthService) LoginUser(userData entities.LoginUserDTO) (string, error) {
	user, err := service.UsersRepository.GetUserByEmail(userData.Email)
	if err != nil {
		return "", err
	}

	if !security.ValidateHashedPassword(userData.Password, user.Password) {
		return "", &customerrors.InvalidPasswordError{}
	}

	// TODO should be changed on JWT
	return "someToken", nil
}

func (service *CommonAuthService) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return service.AuthRepository.RegisterUser(userData)
}
