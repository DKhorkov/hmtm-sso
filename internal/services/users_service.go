package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

func NewCommonUsersService(usersRepository interfaces.UsersRepository, logger *slog.Logger) *CommonUsersService {
	return &CommonUsersService{
		usersRepository: usersRepository,
		logger:          logger,
	}
}

type CommonUsersService struct {
	usersRepository interfaces.UsersRepository
	logger          *slog.Logger
}

func (service *CommonUsersService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return service.usersRepository.GetAllUsers(ctx)
}

func (service *CommonUsersService) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	user, err := service.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get User with ID=%d", id),
			err,
		)

		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}

func (service *CommonUsersService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := service.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get User with Email=%s", email),
			err,
		)

		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}
