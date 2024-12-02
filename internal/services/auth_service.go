package services

import (
	"log/slog"
	"time"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonAuthService struct {
	authRepository  interfaces.AuthRepository
	usersRepository interfaces.UsersRepository
	logger          *slog.Logger
}

func (service *CommonAuthService) RegisterUser(userData entities.RegisterUserDTO) (uint64, error) {
	user, _ := service.usersRepository.GetUserByEmail(userData.Credentials.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.authRepository.RegisterUser(userData)
}

func (service *CommonAuthService) CreateRefreshToken(
	userID uint64,
	refreshToken string,
	ttl time.Duration,
) (uint64, error) {
	return service.authRepository.CreateRefreshToken(
		userID,
		refreshToken,
		ttl,
	)
}

func (service *CommonAuthService) GetRefreshTokenByUserID(userID uint64) (*entities.RefreshToken, error) {
	return service.authRepository.GetRefreshTokenByUserID(userID)
}

func (service *CommonAuthService) ExpireRefreshToken(refreshToken string) error {
	return service.authRepository.ExpireRefreshToken(refreshToken)
}

func NewCommonAuthService(
	authRepository interfaces.AuthRepository,
	usersRepository interfaces.UsersRepository,
	logger *slog.Logger,
) *CommonAuthService {
	return &CommonAuthService{
		authRepository:  authRepository,
		usersRepository: usersRepository,
		logger:          logger,
	}
}
