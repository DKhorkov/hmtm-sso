package repositories

import (
	"time"

	"github.com/DKhorkov/hmtm-sso/internal/entities"

	"github.com/DKhorkov/libs/db"
)

type CommonAuthRepository struct {
	dbConnector db.Connector
}

func (repo *CommonAuthRepository) RegisterUser(userData entities.RegisterUserDTO) (uint64, error) {
	var userID uint64
	connection := repo.dbConnector.GetConnection()
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

func (repo *CommonAuthRepository) CreateRefreshToken(
	userID uint64,
	refreshToken string,
	ttl time.Duration,
) (uint64, error) {
	var refreshTokenID uint64
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func (repo *CommonAuthRepository) GetRefreshTokenByUserID(userID uint64) (*entities.RefreshToken, error) {
	refreshToken := &entities.RefreshToken{}
	columns := db.GetEntityColumns(refreshToken)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func (repo *CommonAuthRepository) ExpireRefreshToken(refreshToken string) error {
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func NewCommonAuthRepository(dbConnector db.Connector) *CommonAuthRepository {
	return &CommonAuthRepository{dbConnector: dbConnector}
}
