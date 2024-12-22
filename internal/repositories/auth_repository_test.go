package repositories_test

import (
	"context"
	"testing"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestRepositoriesRegisterUser(t *testing.T) {
	const testUserEmail = "user@example.com"
	testUserDTO := entities.RegisterUserDTO{
		Credentials: entities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "password",
		},
	}

	t.Run("successful registration", func(t *testing.T) {
		ctx := context.Background()

		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		authRepository := repositories.NewCommonAuthRepository(dbConnector)

		// Error and zero userID due to returning nil ID after register.
		// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
		userID, err := authRepository.RegisterUser(ctx, testUserDTO)
		require.Error(t, err)
		assert.Equal(t, uint64(0), userID)

		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		var usersCount int
		err = connection.QueryRowContext(
			ctx,
			`
				SELECT COUNT(*)
				FROM users
			`,
		).Scan(&usersCount)
		require.NoError(t, err)
		assert.Equal(t, 1, usersCount)
	})

	t.Run("registration failure due to existence of user with same email", func(t *testing.T) {
		ctx := context.Background()

		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		_, err = connection.ExecContext(
			ctx,
			`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
			1,
			testUserDTO.Credentials.Email,
			testUserDTO.Credentials.Password,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		authRepository := repositories.NewCommonAuthRepository(dbConnector)
		userID, err := authRepository.RegisterUser(ctx, testUserDTO)
		require.Error(t, err)
		assert.Equal(t, uint64(0), userID)
	})
}
