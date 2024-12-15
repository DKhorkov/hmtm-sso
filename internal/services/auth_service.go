package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/entities"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

type CommonAuthService struct {
	authRepository  interfaces.AuthRepository
	usersRepository interfaces.UsersRepository
	logger          *slog.Logger
}

func (service *CommonAuthService) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	user, _ := service.usersRepository.GetUserByEmail(userData.Credentials.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.authRepository.RegisterUser(userData)
}

func (service *CommonAuthService) CreateRefreshToken(
	ctx context.Context,
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

func (service *CommonAuthService) GetRefreshTokenByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.RefreshToken, error) {
	return service.authRepository.GetRefreshTokenByUserID(userID)
}

func (service *CommonAuthService) ExpireRefreshToken(ctx context.Context, refreshToken string) error {
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
