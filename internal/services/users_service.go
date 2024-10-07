package services

import (
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"

	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
)

type CommonUsersService struct {
	UsersRepository interfaces.UsersRepository
}

func (service *CommonUsersService) GetAllUsers() ([]*entities.User, error) {
	return service.UsersRepository.GetAllUsers()
}

func (service *CommonUsersService) GetUserByID(id int) (*entities.User, error) {
	user, err := service.UsersRepository.GetUserByID(id)
	if err != nil {
		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}
