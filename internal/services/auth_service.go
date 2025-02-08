package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

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

type CommonAuthService struct {
	authRepository  interfaces.AuthRepository
	usersRepository interfaces.UsersRepository
	logger          *slog.Logger
}

func (service *CommonAuthService) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	user, _ := service.usersRepository.GetUserByEmail(ctx, userData.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.authRepository.RegisterUser(ctx, userData)
}

func (service *CommonAuthService) CreateRefreshToken(
	ctx context.Context,
	userID uint64,
	refreshToken string,
	ttl time.Duration,
) (uint64, error) {
	return service.authRepository.CreateRefreshToken(
		ctx,
		userID,
		refreshToken,
		ttl,
	)
}

func (service *CommonAuthService) GetRefreshTokenByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.RefreshToken, error) {
	return service.authRepository.GetRefreshTokenByUserID(ctx, userID)
}

func (service *CommonAuthService) ExpireRefreshToken(ctx context.Context, refreshToken string) error {
	return service.authRepository.ExpireRefreshToken(ctx, refreshToken)
}

func (service *CommonAuthService) VerifyUserEmail(ctx context.Context, userID uint64) error {
	return service.authRepository.VerifyUserEmail(ctx, userID)
}
