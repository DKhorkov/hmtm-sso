package repositories

import (
	"context"
	"log/slog"
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/libs/db"
)

func NewCommonAuthRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
) *CommonAuthRepository {
	return &CommonAuthRepository{
		dbConnector: dbConnector,
		logger:      logger,
	}
}

type CommonAuthRepository struct {
	dbConnector db.Connector
	logger      *slog.Logger
}

func (repo *CommonAuthRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	var userID uint64
	err = connection.QueryRowContext(
		ctx,
		`
			INSERT INTO users (display_name, email, password) 
			VALUES ($1, $2, $3)
			RETURNING users.id
		`,
		userData.DisplayName,
		userData.Email,
		userData.Password,
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

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

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

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

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

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	_, err = connection.ExecContext(
		ctx,
		`
			UPDATE refresh_tokens
			SET ttl = $1
			WHERE value = $2
		`,
		time.Now().Add(time.Hour*time.Duration(-24)),
		refreshToken,
	)

	return err
}

func (repo *CommonAuthRepository) Close() error {
	return nil
}
