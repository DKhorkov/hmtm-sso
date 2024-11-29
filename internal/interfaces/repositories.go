package interfaces

import (
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type UsersRepository interface {
	GetUserByID(id uint64) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
}

type AuthRepository interface {
	RegisterUser(userData entities.RegisterUserDTO) (userID uint64, err error)
	CreateRefreshToken(userID uint64, refreshToken string, ttl time.Duration) (refreshTokenID uint64, err error)
	GetRefreshTokenByUserID(userID uint64) (*entities.RefreshToken, error)
	ExpireRefreshToken(refreshToken string) error
}
