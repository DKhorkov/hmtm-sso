package interfaces

import (
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type UseCases interface {
	GetUserByID(id int) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	RegisterUser(userData entities.RegisterUserDTO) (userID int, err error)
	LoginUser(userData entities.LoginUserDTO) (*entities.TokensDTO, error)
	GetMe(accessToken string) (*entities.User, error)
	RefreshTokens(refreshTokensData entities.TokensDTO) (*entities.TokensDTO, error)
}
