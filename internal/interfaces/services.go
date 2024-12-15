package interfaces

import (
	"context"
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

type UsersService interface {
	GetUserByID(ctx context.Context, id uint64) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}

type AuthService interface {
	RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (userID uint64, err error)
	CreateRefreshToken(
		ctx context.Context,
		userID uint64,
		refreshToken string,
		ttl time.Duration,
	) (refreshTokenID uint64, err error)
	GetRefreshTokenByUserID(ctx context.Context, userID uint64) (*entities.RefreshToken, error)
	ExpireRefreshToken(ctx context.Context, refreshToken string) error
}
