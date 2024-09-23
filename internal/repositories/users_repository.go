package repositories

import (
	"github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/internal/database"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

type CommonUsersRepository struct {
	DBConnector interfaces.DBConnector
}

func (repo *CommonUsersRepository) GetUserByID(id int) (*entities.User, error) {
	user := &entities.User{}
	columns := database.GetEntityColumns(user)
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT * 
			FROM users AS u
			WHERE u.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}

func (repo *CommonUsersRepository) GetUserByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	columns := database.GetEntityColumns(user)
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT * 
			FROM users AS u
			WHERE u.email = $1
		`,
		email,
	).Scan(columns...)

	if err != nil {
		return nil, &customerrors.UserNotFoundError{}
	}

	return user, nil
}

func (repo *CommonUsersRepository) GetAllUsers() ([]*entities.User, error) {
	connection := repo.DBConnector.GetConnection()
	rows, err := connection.Query(
		`
			SELECT * 
			FROM users
		`,
	)

	if err != nil {
		return nil, err
	}

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		columns := database.GetEntityColumns(user)
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}
