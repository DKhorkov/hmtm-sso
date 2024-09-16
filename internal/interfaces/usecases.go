package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/entities"
)

type UseCases interface {
	GetUserByID(id int) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	RegisterUser(userData entities.RegisterUserDTO) (int, error)
	LoginUser(userData entities.LoginUserDTO) (string, error)
}
