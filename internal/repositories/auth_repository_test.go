//go:build integration

package repositories_test

import (
	"context"
	"log/slog"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/libs/tracing"
	tracingmock "github.com/DKhorkov/libs/tracing/mocks"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
)

var (
	testUserDTO = entities.RegisterUserDTO{
		DisplayName: "test User",
		Email:       testUserEmail,
		Password:    "password",
	}
	logger     = &slog.Logger{}
	spanConfig = tracing.SpanConfig{}
)

func TestRepositoriesRegisterUser(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		dbConnector := StartUp(t)
		traceProvider := tracingmock.NewMockTraceProvider(gomock.NewController(t))
		traceProvider.EXPECT().Span(gomock.Any(), gomock.Any()).Return(
			context.Background(),
			tracingmock.NewMockSpan(),
		).Times(1)

		authRepository := repositories.NewCommonAuthRepository(dbConnector, logger, traceProvider, spanConfig)

		// Error and zero userID due to returning nil ID after register.
		// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
		ctx := context.Background()
		userID, err := authRepository.RegisterUser(ctx, testUserDTO)
		require.Error(t, err)
		assert.Equal(t, uint64(0), userID)

		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		defer func() {
			if err = connection.Close(); err != nil {
				t.Fatal(err)
			}
		}()

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
		dbConnector := StartUp(t)
		ctx := context.Background()
		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		defer func() {
			if err = connection.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		_, err = connection.ExecContext(
			ctx,
			`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
			1,
			testUserDTO.DisplayName,
			testUserDTO.Email,
			testUserDTO.Password,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		traceProvider := tracingmock.NewMockTraceProvider(gomock.NewController(t))
		traceProvider.EXPECT().Span(gomock.Any(), gomock.Any()).Return(
			context.Background(),
			tracingmock.NewMockSpan(),
		).Times(1)

		authRepository := repositories.NewCommonAuthRepository(dbConnector, logger, traceProvider, spanConfig)

		userID, err := authRepository.RegisterUser(ctx, testUserDTO)
		require.Error(t, err)
		assert.Equal(t, uint64(0), userID)
	})
}

func TestRepositoriesVerifyUserEmail(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		dbConnector := StartUp(t)
		ctx := context.Background()
		connection, err := dbConnector.Connection(ctx)
		require.NoError(t, err)

		defer func() {
			if err = connection.Close(); err != nil {
				t.Fatal(err)
			}
		}()

		traceProvider := tracingmock.NewMockTraceProvider(gomock.NewController(t))
		traceProvider.EXPECT().Span(gomock.Any(), gomock.Any()).Return(
			context.Background(),
			tracingmock.NewMockSpan(),
		).Times(1)

		testUser := &entities.User{
			ID:          testUserID,
			DisplayName: "Display Name",
			Email:       testUserEmail,
			Password:    "password",
		}

		_, err = connection.ExecContext(
			ctx,
			`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
			testUser.ID,
			testUser.DisplayName,
			testUser.Email,
			testUser.Password,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		authRepository := repositories.NewCommonAuthRepository(dbConnector, logger, traceProvider, spanConfig)

		// Error and zero userID due to returning nil ID after register.
		// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
		err = authRepository.VerifyUserEmail(ctx, testUserID)
		require.NoError(t, err)

		var emailConfirmed bool
		err = connection.QueryRowContext(
			ctx,
			`
				SELECT u.email_confirmed
				FROM users AS u
				WHERE u.id = $1
			`,
			testUserID,
		).Scan(&emailConfirmed)
		require.NoError(t, err)
		assert.True(t, emailConfirmed)
	})
}

func BenchmarkRepositoriesRegisterUser(b *testing.B) {
	dbConnector := StartUp(b)
	traceProvider := tracingmock.NewMockTraceProvider(gomock.NewController(b))
	traceProvider.EXPECT().Span(gomock.Any(), gomock.Any()).Return(
		context.Background(),
		tracingmock.NewMockSpan(),
	).AnyTimes()

	authRepository := repositories.NewCommonAuthRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = authRepository.RegisterUser(
			context.Background(),
			testUserDTO,
		)
	}
}
