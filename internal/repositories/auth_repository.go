package repositories

import (
	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

type CommonAuthRepository struct {
}

func (repo *CommonAuthRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return 0, nil
}
