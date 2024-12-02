package repositories_test

import (
	"testing"

	"github.com/DKhorkov/hmtm-sso/internal/repositories"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestRepositoriesGetUserByID(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "user@example.com"
	)

	t.Run("get existing user", func(t *testing.T) {
		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		testUser := &entities.User{
			ID:       testUserID,
			Email:    testUserEmail,
			Password: "password",
		}

		_, err := dbConnector.GetConnection().Exec(
			`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
			testUser.ID,
			testUser.Email,
			testUser.Password,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		usersRepository := repositories.NewCommonUsersRepository(dbConnector)
		user, err := usersRepository.GetUserByID(testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("get non existing user failure", func(t *testing.T) {
		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		usersRepository := repositories.NewCommonUsersRepository(dbConnector)
		user, err := usersRepository.GetUserByID(testUserID)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestRepositoriesGetUserByEmail(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "user@example.com"
	)

	t.Run("get existing user", func(t *testing.T) {
		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		testUser := &entities.User{
			ID:       testUserID,
			Email:    testUserEmail,
			Password: "password",
		}

		_, err := dbConnector.GetConnection().Exec(
			`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
			testUser.ID,
			testUser.Email,
			testUser.Password,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		usersRepository := repositories.NewCommonUsersRepository(dbConnector)
		user, err := usersRepository.GetUserByEmail(testUser.Email)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("get non existing user failure", func(t *testing.T) {
		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		usersRepository := repositories.NewCommonUsersRepository(dbConnector)
		user, err := usersRepository.GetUserByEmail(testUserEmail)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestRepositoriesGetAllUsers(t *testing.T) {
	t.Run("get users with existing users", func(t *testing.T) {
		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		testUser := &entities.User{
			ID:       1,
			Email:    "user@example.com",
			Password: "password",
		}

		_, err := dbConnector.GetConnection().Exec(
			`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
			testUser.ID,
			testUser.Email,
			testUser.Password,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		usersRepository := repositories.NewCommonUsersRepository(dbConnector)
		users, err := usersRepository.GetAllUsers()
		require.NoError(t, err)
		assert.IsType(t, []*entities.User{}, users)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 1)
	})

	t.Run("get users without existing users", func(t *testing.T) {
		dbConnector := StartUp(t)
		defer TearDown(t, dbConnector)

		usersRepository := repositories.NewCommonUsersRepository(dbConnector)
		users, err := usersRepository.GetAllUsers()
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}
