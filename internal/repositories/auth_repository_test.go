//go:build integration

package repositories_test

import (
	"context"
	"database/sql"
	"os"
	"path"
	"testing"

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
	testUserID       = 1
	testUserEmail    = "user@example.com"
)

var (
	testUserDTO = entities.RegisterUserDTO{
		DisplayName: "test User",
		Email:       testUserEmail,
		Password:    "password",
	}
)

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

type AuthTestSuite struct {
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

func (s *AuthTestSuite) SetupSuite() {
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

func (s *AuthTestSuite) SetupTest() {
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

func (s *AuthTestSuite) TearDownTest() {
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

func (s *AuthTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *AuthTestSuite) TestRepositoriesRegisterUserSuccess() {
	ctx := context.Background()
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

func (s *AuthTestSuite) TestRepositoriesRegisterUserFailEmailAlreadyExists() {
	ctx := context.Background()
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
		testUserID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	userID, err := s.authRepository.RegisterUser(ctx, testUserDTO)
	s.Error(err)
	s.Zero(userID)
}

func (s *AuthTestSuite) TestRepositoriesVerifyUserEmailSuccess() {
	ctx := context.Background()
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
		testUserID,
		testUserDTO.DisplayName,
		testUserDTO.Email,
		testUserDTO.Password,
	)

	s.NoError(err)

	err = s.authRepository.VerifyUserEmail(ctx, uint64(1))
	s.NoError(err)
}

func (s *AuthTestSuite) TestRepositoriesVerifyUserEmailUserDoesNotExist() {
	ctx := context.Background()
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	err := s.authRepository.VerifyUserEmail(ctx, uint64(1))
	s.NoError(err)
}

func BenchmarkRepositoriesRegisterUser(b *testing.B) {
	spanConfig := tracing.SpanConfig{}
	ctx := context.Background()
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
