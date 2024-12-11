package services

import (
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/internal/entities"

	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/libs/logging"
)

type CommonUsersService struct {
	usersRepository interfaces.UsersRepository
	logger          *slog.Logger
}

func (service *CommonUsersService) GetAllUsers() ([]entities.User, error) {
	return service.usersRepository.GetAllUsers()
}

func (service *CommonUsersService) GetUserByID(id uint64) (*entities.User, error) {
	user, err := service.usersRepository.GetUserByID(id)
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get user by id",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customerrors.UserNotFoundError{BaseErr: err}
	}

	return user, nil
}

func (service *CommonUsersService) GetUserByEmail(email string) (*entities.User, error) {
	user, err := service.usersRepository.GetUserByEmail(email)
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get user by email",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customerrors.UserNotFoundError{BaseErr: err}
	}

	return user, nil
}

func NewCommonUsersService(usersRepository interfaces.UsersRepository, logger *slog.Logger) *CommonUsersService {
	return &CommonUsersService{
		usersRepository: usersRepository,
		logger:          logger,
	}
}
