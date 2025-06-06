package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

//go:generate mockgen -source=usecases.go -destination=../../mocks/usecases/usecases.go -package=mockusecases
type UseCases interface {
	GetUserByID(ctx context.Context, id uint64) (*entities.User, error)
	GetUsers(ctx context.Context, pagination *entities.Pagination) ([]entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetMe(ctx context.Context, accessToken string) (*entities.User, error)
	UpdateUserProfile(
		ctx context.Context,
		rawUserProfileData entities.RawUpdateUserProfileDTO,
	) error

	RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (userID uint64, err error)
	LoginUser(ctx context.Context, userData entities.LoginUserDTO) (*entities.TokensDTO, error)
	LogoutUser(ctx context.Context, accessToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error)
	VerifyUserEmail(ctx context.Context, verifyEmailToken string) error
	ForgetPassword(ctx context.Context, forgetPasswordToken, newPassword string) error
	SendForgetPasswordMessage(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, accessToken, oldPassword, newPassword string) error
	SendVerifyEmailMessage(ctx context.Context, email string) error
}
