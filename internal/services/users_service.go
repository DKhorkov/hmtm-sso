package services

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

type CommonUsersService struct {
	UsersRepository interfaces.UsersRepository
}

func (service *CommonUsersService) GetAllUsers() ([]*entities.User, error) {
	return service.UsersRepository.GetAllUsers()
}

func (service *CommonUsersService) GetUserByID(id int) (*entities.User, error) {
	return service.UsersRepository.GetUserByID(id)
}
