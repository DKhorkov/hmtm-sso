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
	RegisterUser(userData entities.RegisterUserDTO) (userID int, err error)
	CreateRefreshToken(userID int, refreshToken string, ttl time.Duration) (refreshTokenID int, err error)
	GetRefreshTokenByUserID(userID int) (*entities.RefreshToken, error)
	ExpireRefreshToken(refreshToken string) error
}
