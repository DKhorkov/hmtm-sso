package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/entities"
)

type UsersService interface {
	GetAllUsers() ([]*entities.User, error)
	GetUserByID(int) (*entities.User, error)
}

type AuthService interface {
	LoginUser(userData entities.LoginUserDTO) (string, error)
	RegisterUser(userData entities.RegisterUserDTO) (int, error)
}
