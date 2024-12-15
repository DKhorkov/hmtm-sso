package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

type UseCases interface {
	GetUserByID(ctx context.Context, id uint64) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (userID uint64, err error)
	LoginUser(ctx context.Context, userData entities.LoginUserDTO) (*entities.TokensDTO, error)
	GetMe(ctx context.Context, accessToken string) (*entities.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error)
}
