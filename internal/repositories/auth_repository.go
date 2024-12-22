package repositories

import (
	"context"
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/libs/db"
)

func NewCommonAuthRepository(dbConnector db.Connector) *CommonAuthRepository {
	return &CommonAuthRepository{dbConnector: dbConnector}
}

type CommonAuthRepository struct {
	dbConnector db.Connector
}

func (repo *CommonAuthRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	var userID uint64
	err = connection.QueryRowContext(
		ctx,
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

func (repo *CommonAuthRepository) CreateRefreshToken(
	ctx context.Context,
	userID uint64,
	refreshToken string,
	ttl time.Duration,
) (uint64, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	var refreshTokenID uint64
	err = connection.QueryRowContext(
		ctx,
		`
			INSERT INTO refresh_tokens (user_id, value, ttl) 
			VALUES ($1, $2, $3)
			RETURNING refresh_tokens.id
		`,
		userID,
		refreshToken,
		time.Now().Add(ttl),
	).Scan(&refreshTokenID)

	if err != nil {
		return 0, err
	}

	return refreshTokenID, nil
}

func (repo *CommonAuthRepository) GetRefreshTokenByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.RefreshToken, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	refreshToken := &entities.RefreshToken{}
	columns := db.GetEntityColumns(refreshToken)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM refresh_tokens AS rt
			WHERE rt.user_id = $1
			  AND rt.ttl > CURRENT_TIMESTAMP
		`,
		userID,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (repo *CommonAuthRepository) ExpireRefreshToken(ctx context.Context, refreshToken string) error {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	err = connection.QueryRowContext(
		ctx,
		`
			UPDATE refresh_tokens
			SET ttl = $1
			WHERE value = $2
		`,
		time.Now().Add(time.Hour*time.Duration(-24)),
		refreshToken,
	).Err()

	return err
}

func (repo *CommonAuthRepository) Close() error {
	return nil
}
