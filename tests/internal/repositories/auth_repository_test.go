package repositories__test

import (
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	"testing"

	"github.com/DKhorkov/hmtm-sso/internal/database"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
	testlifespan "github.com/DKhorkov/hmtm-sso/tests/internal/repositories/lifespan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestRegisterUser(t *testing.T) {
	const testUserEmail = "user@example.com"
	testUserDTO := entities.RegisterUserDTO{
		Credentials: entities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "password",
		},
	}

	t.Run("successful registration", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		// Error and zero userID due to returning nil ID after register.
		// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
		userID, err := authRepository.RegisterUser(testUserDTO)
		require.Error(t, err)
		assert.Equal(t, 0, userID)

		var usersCount int
		err = connection.QueryRow(
			`
				SELECT COUNT(*)
				FROM users
			`,
		).Scan(&usersCount)
		require.NoError(t, err)
		assert.Equal(t, 1, usersCount)
	})

	t.Run("registration failure due to existence of user with same email", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		_, err := connection.Exec(
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

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		userID, err := authRepository.RegisterUser(testUserDTO)
		require.Error(t, err)
		assert.Equal(t, 0, userID)
	})
}
