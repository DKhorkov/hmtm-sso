package services

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/internal/config"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/security"
)

type CommonAuthService struct {
	AuthRepository  interfaces.AuthRepository
	UsersRepository interfaces.UsersRepository
	JWTConfig       config.JWTConfig
}

func (service *CommonAuthService) LoginUser(userData entities.LoginUserDTO) (string, error) {
	user, err := service.UsersRepository.GetUserByEmail(userData.Email)
	if err != nil {
		return "", err
	}

	if !security.ValidateHashedPassword(userData.Password, user.Password) {
		return "", &customerrors.InvalidPasswordError{}
	}

	return security.GenerateJWT(
		user,
		service.JWTConfig.SecretKey,
		service.JWTConfig.TTL,
		service.JWTConfig.Algorithm,
	)
}

func (service *CommonAuthService) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return service.AuthRepository.RegisterUser(userData)
}
