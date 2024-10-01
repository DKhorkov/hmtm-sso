package services

import (
	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/security"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
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
	user, _ := service.UsersRepository.GetUserByEmail(userData.Credentials.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.AuthRepository.RegisterUser(userData)
}
