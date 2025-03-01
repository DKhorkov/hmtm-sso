//go:build integration

package repositories_test

import (
	"context"
	"database/sql"
	"os"
	"path"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/libs/db"
	loggermock "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/tracing"
	tracingmock "github.com/DKhorkov/libs/tracing/mocks"

	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
)

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

type UsersTestSuite struct {
	suite.Suite

	cwd             string
	ctx             context.Context
	dbConnector     db.Connector
	connection      *sql.Conn
	usersRepository interfaces.UsersRepository
	logger          *loggermock.MockLogger
	traceProvider   *tracingmock.MockProvider
	spanConfig      tracing.SpanConfig
}

func (s *UsersTestSuite) SetupSuite() {
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
	s.usersRepository = repositories.NewUsersRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *UsersTestSuite) SetupTest() {
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

func (s *UsersTestSuite) TearDownTest() {
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

func (s *UsersTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *UsersTestSuite) TestRepositoriesGetExistingUserByID() {
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

	user, err := s.usersRepository.GetUserByID(ctx, testUserID)
	s.NoError(err)
	s.NotNil(user)
}

func (s *UsersTestSuite) TestRepositoriesGetNonExistingUserByID() {
	ctx := context.Background()
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	user, err := s.usersRepository.GetUserByID(ctx, testUserID)
	s.Error(err)
	s.Nil(user)
}

func (s *UsersTestSuite) TestRepositoriesGetExistingUserByEmail() {
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

	user, err := s.usersRepository.GetUserByEmail(ctx, testUserEmail)
	s.NoError(err)
	s.NotNil(user)
}

func (s *UsersTestSuite) TestRepositoriesGetNonExistingUserByEmail() {
	ctx := context.Background()
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	user, err := s.usersRepository.GetUserByEmail(ctx, testUserEmail)
	s.Error(err)
	s.Nil(user)
}

func (s *UsersTestSuite) TestRepositoriesGetAllUserWithExistingUsers() {
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

	users, err := s.usersRepository.GetAllUsers(ctx)
	s.NoError(err)
	s.NotEmpty(users)
}

func (s *UsersTestSuite) TestRepositoriesGetAllUserWithoutExistingUsers() {
	ctx := context.Background()
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), tracingmock.NewMockSpan()).
		Times(1)

	users, err := s.usersRepository.GetAllUsers(ctx)
	s.NoError(err)
	s.Empty(users)
}

func BenchmarkRepositoriesGetUserByID(b *testing.B) {
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

	usersRepository := repositories.NewUsersRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = usersRepository.GetUserByID(
			ctx,
			testUserID,
		)
	}
}

func BenchmarkRepositoriesGetUserByEmail(b *testing.B) {
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

	usersRepository := repositories.NewUsersRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = usersRepository.GetUserByEmail(
			ctx,
			testUserEmail,
		)
	}
}

func BenchmarkRepositoriesGetAllUsers(b *testing.B) {
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

	usersRepository := repositories.NewUsersRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = usersRepository.GetAllUsers(ctx)
	}
}
