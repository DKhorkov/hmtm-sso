package repositories

import (
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonAuthRepository struct {
	DBConnector interfaces.DBConnector
}

func (repo *CommonAuthRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	var userID int
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			INSERT INTO users (email, password) 
			VALUES ($1, $2)
			RETURNING users.id
		`,
		userData.Credentials.Email,
		userData.Credentials.Password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}
