package services

import (
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonUsersService struct {
	usersRepository interfaces.UsersRepository
}

func (service *CommonUsersService) GetAllUsers() ([]*entities.User, error) {
	return service.usersRepository.GetAllUsers()
}

func (service *CommonUsersService) GetUserByID(id int) (*entities.User, error) {
	user, err := service.usersRepository.GetUserByID(id)
	if err != nil {
		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}

func (service *CommonUsersService) GetUserByEmail(email string) (*entities.User, error) {
	user, err := service.usersRepository.GetUserByEmail(email)
	if err != nil {
		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}

func NewCommonUsersService(usersRepository interfaces.UsersRepository) *CommonUsersService {
	return &CommonUsersService{usersRepository: usersRepository}
}
