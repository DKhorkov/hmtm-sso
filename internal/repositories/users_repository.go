package repositories

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
)

type CommonUsersRepository struct {
}

func (repo *CommonUsersRepository) GetUserByID(id int) (*entities.User, error) {
	return nil, &customerrors.UserNotFoundError{}
}

func (repo *CommonUsersRepository) GetUserByEmail(email string) (*entities.User, error) {
	return nil, &customerrors.UserNotFoundError{}
}

func (repo *CommonUsersRepository) GetAllUsers() ([]*entities.User, error) {
	return []*entities.User{}, nil
}
