package interfaces

import (
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type UsersService interface {
	GetAllUsers() ([]*entities.User, error)
	GetUserByID(id int) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
}

type AuthService interface {
	RegisterUser(userData entities.RegisterUserDTO) (int, error)
	CreateRefreshToken(userID int, refreshToken string, ttl time.Duration) (int, error)
	GetRefreshTokenByUserID(userID int) (*entities.RefreshToken, error)
	ExpireRefreshToken(refreshToken string) error
}
