package services

import (
	"time"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonAuthService struct {
	AuthRepository  interfaces.AuthRepository
	UsersRepository interfaces.UsersRepository
}

func (service *CommonAuthService) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	user, _ := service.UsersRepository.GetUserByEmail(userData.Credentials.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.AuthRepository.RegisterUser(userData)
}

func (service *CommonAuthService) CreateRefreshToken(
	userID int,
	refreshToken string,
	ttl time.Duration,
) (int, error) {
	return service.AuthRepository.CreateRefreshToken(userID, refreshToken, ttl)
}

func (service *CommonAuthService) GetRefreshTokenByUserID(userID int) (*entities.RefreshToken, error) {
	return service.AuthRepository.GetRefreshTokenByUserID(userID)
}

func (service *CommonAuthService) ExpireRefreshToken(refreshToken string) error {
	return service.AuthRepository.ExpireRefreshToken(refreshToken)
}
