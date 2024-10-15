package interfaces

import (
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type UsersRepository interface {
	GetUserByID(id int) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
}

type AuthRepository interface {
	RegisterUser(user entities.RegisterUserDTO) (int, error)
	CreateRefreshToken(userID int, refreshToken string, ttl time.Duration) (int, error)
	GetRefreshTokenByUserID(userID int) (*entities.RefreshToken, error)
	ExpireRefreshToken(refreshToken string) error
}
