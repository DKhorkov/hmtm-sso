package services

import (
	"context"
	"time"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

type AuthService struct {
	authRepository  interfaces.AuthRepository
	usersRepository interfaces.UsersRepository
	logger          logging.Logger
}

func NewAuthService(
	authRepository interfaces.AuthRepository,
	usersRepository interfaces.UsersRepository,
	logger logging.Logger,
) *AuthService {
	return &AuthService{
		authRepository:  authRepository,
		usersRepository: usersRepository,
		logger:          logger,
	}
}

func (service *AuthService) RegisterUser(
	ctx context.Context,
	userData entities.RegisterUserDTO,
) (uint64, error) {
	user, _ := service.usersRepository.GetUserByEmail(ctx, userData.Email)
	if user != nil {
		return 0, &customerrors.UserAlreadyExistsError{}
	}

	return service.authRepository.RegisterUser(ctx, userData)
}

func (service *AuthService) CreateRefreshToken(
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

func (service *AuthService) GetRefreshTokenByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.RefreshToken, error) {
	return service.authRepository.GetRefreshTokenByUserID(ctx, userID)
}

func (service *AuthService) ExpireRefreshToken(ctx context.Context, refreshToken string) error {
	return service.authRepository.ExpireRefreshToken(ctx, refreshToken)
}

func (service *AuthService) VerifyUserEmail(ctx context.Context, userID uint64) error {
	return service.authRepository.VerifyUserEmail(ctx, userID)
}

func (service *AuthService) ForgetPassword(
	ctx context.Context,
	userID uint64,
	newPassword string,
) error {
	return service.authRepository.ForgetPassword(ctx, userID, newPassword)
}

func (service *AuthService) ChangePassword(
	ctx context.Context,
	userID uint64,
	newPassword string,
) error {
	return service.authRepository.ChangePassword(ctx, userID, newPassword)
}
