package repositories

import (
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	"github.com/DKhorkov/libs/db"
)

type CommonUsersRepository struct {
	dbConnector db.Connector
}

func (repo *CommonUsersRepository) GetUserByID(id uint64) (*entities.User, error) {
	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT * 
			FROM users AS u
			WHERE u.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *CommonUsersRepository) GetUserByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT * 
			FROM users AS u
			WHERE u.email = $1
		`,
		email,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *CommonUsersRepository) GetAllUsers() ([]*entities.User, error) {
	connection := repo.dbConnector.GetConnection()
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
		columns := db.GetEntityColumns(user)
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

func NewCommonUsersRepository(dbConnector db.Connector) *CommonUsersRepository {
	return &CommonUsersRepository{dbConnector: dbConnector}
}
