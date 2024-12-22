package repositories

import (
	"context"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/libs/db"
)

func NewCommonUsersRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
) *CommonUsersRepository {
	return &CommonUsersRepository{
		dbConnector: dbConnector,
		logger:      logger,
	}
}

type CommonUsersRepository struct {
	dbConnector db.Connector
	logger      *slog.Logger
}

func (repo *CommonUsersRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	err = connection.QueryRowContext(
		ctx,
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

func (repo *CommonUsersRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	user := &entities.User{}
	columns := db.GetEntityColumns(user)
	err = connection.QueryRowContext(
		ctx,
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

func (repo *CommonUsersRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM users
		`,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			logging.LogErrorContext(
				ctx,
				repo.logger,
				"error during closing SQL rows",
				err,
			)
		}
	}()

	var users []entities.User
	for rows.Next() {
		user := entities.User{}
		columns := db.GetEntityColumns(&user) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *CommonUsersRepository) Close() error {
	return nil
}
