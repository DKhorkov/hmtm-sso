package repositories_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/DKhorkov/hmtm-sso/internal/entities"

	"github.com/DKhorkov/hmtm-sso/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

const (
	testUserID    = 1
	testUserEmail = "user@example.com"
)

func TestRepositoriesGetUserByID(t *testing.T) {
	t.Run("get existing user", func(t *testing.T) {
		dbConnector := StartUp(t)
		ctx := context.Background()
		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		testUser := &entities.User{
			ID:       testUserID,
			Email:    testUserEmail,
			Password: "password",
		}

		_, err = connection.ExecContext(
			ctx,
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

		usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

		user, err := usersRepository.GetUserByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("get non existing user failure", func(t *testing.T) {
		dbConnector := StartUp(t)
		usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

		user, err := usersRepository.GetUserByID(context.Background(), testUserID)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestRepositoriesGetUserByEmail(t *testing.T) {
	t.Run("get existing user", func(t *testing.T) {
		dbConnector := StartUp(t)
		ctx := context.Background()
		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		testUser := &entities.User{
			ID:       testUserID,
			Email:    testUserEmail,
			Password: "password",
		}

		_, err = connection.ExecContext(
			ctx,
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

		usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

		user, err := usersRepository.GetUserByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("get non existing user failure", func(t *testing.T) {
		dbConnector := StartUp(t)
		usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

		user, err := usersRepository.GetUserByEmail(context.Background(), testUserEmail)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestRepositoriesGetAllUsers(t *testing.T) {
	t.Run("get users with existing users", func(t *testing.T) {
		dbConnector := StartUp(t)
		ctx := context.Background()
		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		testUser := &entities.User{
			ID:       1,
			Email:    "user@example.com",
			Password: "password",
		}

		_, err = connection.ExecContext(
			ctx,
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

		usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

		users, err := usersRepository.GetAllUsers(ctx)
		require.NoError(t, err)
		assert.IsType(t, []entities.User{}, users)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 1)
	})

	t.Run("get users without existing users", func(t *testing.T) {
		dbConnector := StartUp(t)
		usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

		users, err := usersRepository.GetAllUsers(context.Background())
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}

func BenchmarkRepositoriesGetUserByID(b *testing.B) {
	dbConnector := StartUp(b)
	ctx := context.Background()
	connection, err := dbConnector.Connection(ctx)
	require.NoError(b, err)
	_, err = connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
		testUserID,
		testUserEmail,
		"testUserPassword",
	)

	if err != nil {
		b.Fatalf("failed to insert user: %v", err)
	}

	usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = usersRepository.GetUserByID(
			ctx,
			testUserID,
		)
	}
}

func BenchmarkRepositoriesGetUserByEmail(b *testing.B) {
	dbConnector := StartUp(b)
	ctx := context.Background()
	connection, err := dbConnector.Connection(ctx)
	require.NoError(b, err)
	_, err = connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
		testUserID,
		testUserEmail,
		"testUserPassword",
	)

	if err != nil {
		b.Fatalf("failed to insert user: %v", err)
	}

	usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = usersRepository.GetUserByEmail(
			ctx,
			testUserEmail,
		)
	}
}

func BenchmarkRepositoriesGetAllUsers(b *testing.B) {
	dbConnector := StartUp(b)
	ctx := context.Background()
	connection, err := dbConnector.Connection(ctx)
	require.NoError(b, err)
	_, err = connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, email, password) 
				VALUES ($1, $2, $3)
			`,
		testUserID,
		testUserEmail,
		"testUserPassword",
	)

	if err != nil {
		b.Fatalf("failed to insert user: %v", err)
	}

	usersRepository := repositories.NewCommonUsersRepository(dbConnector, &slog.Logger{})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = usersRepository.GetAllUsers(ctx)
	}
}
