package services

import (
	"time"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonAuthService struct {
	authRepository  interfaces.AuthRepository
	usersRepository interfaces.UsersRepository
}

func (service *CommonAuthService) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	user, _ := service.usersRepository.GetUserByEmail(userData.Credentials.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.authRepository.RegisterUser(userData)
}

func (service *CommonAuthService) CreateRefreshToken(
	userID int,
	refreshToken string,
	ttl time.Duration,
) (int, error) {
	return service.authRepository.CreateRefreshToken(userID, refreshToken, ttl)
}

func (service *CommonAuthService) GetRefreshTokenByUserID(userID int) (*entities.RefreshToken, error) {
	return service.authRepository.GetRefreshTokenByUserID(userID)
}

func (service *CommonAuthService) ExpireRefreshToken(refreshToken string) error {
	return service.authRepository.ExpireRefreshToken(refreshToken)
}

func NewCommonAuthService(
	authRepository interfaces.AuthRepository,
	usersRepository interfaces.UsersRepository,
) *CommonAuthService {
	return &CommonAuthService{
		authRepository:  authRepository,
		usersRepository: usersRepository,
	}
}
