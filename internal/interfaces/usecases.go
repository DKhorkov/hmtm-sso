package interfaces

import "github.com/DKhorkov/hmtm-sso/internal/entities"

type UseCases interface {
	GetUserByID(id uint64) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
	RegisterUser(userData entities.RegisterUserDTO) (userID uint64, err error)
	LoginUser(userData entities.LoginUserDTO) (*entities.TokensDTO, error)
	GetMe(accessToken string) (*entities.User, error)
	RefreshTokens(refreshToken string) (*entities.TokensDTO, error)
}
