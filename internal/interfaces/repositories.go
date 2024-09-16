package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/entities"
)

type UsersRepository interface {
	GetUserByID(id int) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
}

type AuthRepository interface {
	RegisterUser(user entities.RegisterUserDTO) (int, error)
}
