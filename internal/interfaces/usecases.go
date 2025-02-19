package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

type UseCases interface {
	UsersService
	GetMe(ctx context.Context, accessToken string) (*entities.User, error)

	RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (userID uint64, err error)
	LoginUser(ctx context.Context, userData entities.LoginUserDTO) (*entities.TokensDTO, error)
	LogoutUser(ctx context.Context, accessToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error)
	VerifyUserEmail(ctx context.Context, verifyEmailToken string) error
}
