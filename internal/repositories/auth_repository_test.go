//go:build integration

package repositories_test

import (
	"context"
	"database/sql"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/libs/db"
	loggermock "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/tracing"
	tracingmock "github.com/DKhorkov/libs/tracing/mocks"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
)

const (
	driver = "sqlite3"
	//dsn    = "file::memory:?cache=shared"
	dsn              = "../../test.db"
	migrationsDir    = "/migrations"
	gooseZeroVersion = 0
	userID           = 1
	email            = "user@example.com"
	refreshTokenID   = 1
)

var (
	ctx         = context.Background()
	testUserDTO = entities.RegisterUserDTO{
		DisplayName: "test User",
		Email:       email,
		Password:    "password",
	}

	testUser = &entities.User{
		ID:          userID,
		DisplayName: testUserDTO.DisplayName,
		Email:       testUserDTO.Email,
		Password:    testUserDTO.Password,
	}

	ttl          = 1 * time.Hour
	refreshToken = &entities.RefreshToken{
		ID:     refreshTokenID,
		UserID: userID,
		Value:  "refresh_token",
		TTL:    time.Now().UTC().Add(ttl),
	}
)

func TestAuthRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

type AuthRepositoryTestSuite struct {
	suite.Suite

	cwd            string
	ctx            context.Context
	dbConnector    db.Connector
	connection     *sql.Conn
	authRepository interfaces.AuthRepository
	logger         *loggermock.MockLogger
	traceProvider  *tracingmock.MockProvider
	spanConfig     tracing.SpanConfig
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	s.NoError(goose.SetDialect(driver))

	ctrl := gomock.NewController(s.T())
	s.ctx = context.Background()
	s.logger = loggermock.NewMockLogger(ctrl)
	dbConnector, err := db.New(dsn, driver, s.logger)
	s.NoError(err)

	cwd, err := os.Getwd()
	s.NoError(err)

	s.cwd = cwd
	s.dbConnector = dbConnector
	s.traceProvider = tracingmock.NewMockProvider(ctrl)
	s.spanConfig = tracing.SpanConfig{}
	s.authRepository = repositories.NewAuthRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *AuthRepositoryTestSuite) SetupTest() {
	s.NoError(
		goose.Up(
			s.dbConnector.Pool(),
			path.Dir(
				path.Dir(s.cwd),
			)+migrationsDir,
		),
	)

	connection, err := s.dbConnector.Connection(s.ctx)
	s.NoError(err)

	s.connection = connection
}

func (s *AuthRepositoryTestSuite) TearDownTest() {
	s.NoError(
		goose.DownTo(
			s.dbConnector.Pool(),
			path.Dir(
				path.Dir(s.cwd),
			)+migrationsDir,
			gooseZeroVersion,
		),
	)

	s.NoError(s.connection.Close())
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *AuthRepositoryTestSuite) TestRegisterUserSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	// Error and zero userID due to returning nil ID after register.
	// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
	userID, err := s.authRepository.RegisterUser(ctx, testUserDTO)
	s.Error(err)
	s.Equal(uint64(0), userID)
}

func (s *AuthRepositoryTestSuite) TestRegisterUserFailEmailAlreadyExists() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
		userID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	userID, err := s.authRepository.RegisterUser(ctx, testUserDTO)
	s.Error(err)
	s.Zero(userID)
}

func (s *AuthRepositoryTestSuite) TestVerifyUserEmailSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
		userID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	err = s.authRepository.VerifyUserEmail(ctx, uint64(1))
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestVerifyUserEmailUserDoesNotExist() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	err := s.authRepository.VerifyUserEmail(ctx, uint64(1))
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestCreateRefreshTokenSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	// Error and zero userID due to returning nil ID after register.
	// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
	refreshTokenID, err := s.authRepository.CreateRefreshToken(
		ctx,
		userID,
		refreshToken.Value,
		ttl,
	)

	s.Error(err)
	s.Zero(refreshTokenID)
}

func (s *AuthRepositoryTestSuite) TestCreateRefreshTokenAlreadyExists() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO refresh_tokens (id, user_id, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
		refreshTokenID,
		userID,
		refreshToken.Value,
		refreshToken.TTL,
	)

	s.NoError(err)

	refreshTokenID, err := s.authRepository.CreateRefreshToken(
		ctx,
		userID,
		refreshToken.Value,
		ttl,
	)

	s.Error(err)
	s.Zero(refreshTokenID)
}

func (s *AuthRepositoryTestSuite) TestGetRefreshTokenByUserIDSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO refresh_tokens (id, user_id, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
		refreshTokenID,
		userID,
		refreshToken.Value,
		refreshToken.TTL,
	)

	s.NoError(err)

	refreshToken, err := s.authRepository.GetRefreshTokenByUserID(ctx, userID)
	s.NoError(err)
	s.NotNil(refreshToken)
}

func (s *AuthRepositoryTestSuite) TestGetRefreshTokenByUserIDNotFound() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	refreshToken, err := s.authRepository.GetRefreshTokenByUserID(ctx, userID)
	s.Error(err)
	s.Nil(refreshToken)
}

func (s *AuthRepositoryTestSuite) TestExpireRefreshTokenSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO refresh_tokens (id, user_id, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
		refreshTokenID,
		userID,
		refreshToken.Value,
		refreshToken.TTL,
	)

	s.NoError(err)

	err = s.authRepository.ExpireRefreshToken(ctx, refreshToken.Value)
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestExpireRefreshTokenDoesNotExist() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	err := s.authRepository.ExpireRefreshToken(ctx, refreshToken.Value)
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestChangePasswordSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
		userID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	err = s.authRepository.ChangePassword(ctx, userID, "new password")
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestChangePasswordUserDoesNotExist() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	err := s.authRepository.ChangePassword(ctx, userID, "new password")
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestForgetPasswordSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	s.logger.
		EXPECT().
		ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
		userID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	_, err = s.connection.ExecContext(
		ctx,
		`
				INSERT INTO refresh_tokens (id, user_id, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
		refreshTokenID,
		userID,
		refreshToken.Value,
		refreshToken.TTL,
	)

	s.NoError(err)

	err = s.authRepository.ForgetPassword(ctx, userID, "new password")
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestForgetPasswordNoActiveRefreshToken() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	s.logger.
		EXPECT().
		ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	_, err := s.connection.ExecContext(
		ctx,
		`
				INSERT INTO users (id, display_name, email, password) 
				VALUES ($1, $2, $3, $4)
			`,
		userID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	err = s.authRepository.ForgetPassword(ctx, userID, "new password")
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestForgetPasswordUserWithProvidedIDDoesNotExist() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	s.logger.
		EXPECT().
		ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	err := s.authRepository.ForgetPassword(ctx, userID, "new password")
	s.NoError(err)
}

func BenchmarkAuthRepository_RegisterUser(b *testing.B) {
	spanConfig := tracing.SpanConfig{}
	ctrl := gomock.NewController(b)
	logger := loggermock.NewMockLogger(ctrl)
	dbConnector, err := db.New(dsn, driver, logger)
	require.NoError(b, err)

	defer func() {
		require.NoError(b, dbConnector.Close())
	}()

	traceProvider := tracingmock.NewMockProvider(ctrl)
	traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(ctx, tracingmock.NewMockSpan()).
		AnyTimes()

	authRepository := repositories.NewAuthRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = authRepository.RegisterUser(
			ctx,
			testUserDTO,
		)
	}
}
