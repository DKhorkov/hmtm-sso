//go:build integration

package repositories_test

import (
	"context"
	"database/sql"
	"github.com/DKhorkov/libs/pointers"
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
	mocktracing "github.com/DKhorkov/libs/tracing/mocks"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/internal/repositories"
)

func TestUsersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UsersRepositoryTestSuite))
}

type UsersRepositoryTestSuite struct {
	suite.Suite

	cwd             string
	ctx             context.Context
	dbConnector     db.Connector
	connection      *sql.Conn
	usersRepository interfaces.UsersRepository
	logger          *loggermock.MockLogger
	traceProvider   *mocktracing.MockProvider
	spanConfig      tracing.SpanConfig
}

func (s *UsersRepositoryTestSuite) SetupSuite() {
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
	s.traceProvider = mocktracing.NewMockProvider(ctrl)
	s.spanConfig = tracing.SpanConfig{}
	s.usersRepository = repositories.NewUsersRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *UsersRepositoryTestSuite) SetupTest() {
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

func (s *UsersRepositoryTestSuite) TearDownTest() {
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

func (s *UsersRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *UsersRepositoryTestSuite) TestGetExistingUserByID() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
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

	user, err := s.usersRepository.GetUserByID(ctx, userID)
	s.NoError(err)
	s.NotNil(user)
}

func (s *UsersRepositoryTestSuite) TestGetNonExistingUserByID() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	user, err := s.usersRepository.GetUserByID(ctx, userID)
	s.Error(err)
	s.Nil(user)
}

func (s *UsersRepositoryTestSuite) TestGetExistingUserByEmail() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
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

	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	s.NoError(err)
	s.NotNil(user)
}

func (s *UsersRepositoryTestSuite) TestGetNonExistingUserByEmail() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	s.Error(err)
	s.Nil(user)
}

func (s *UsersRepositoryTestSuite) TestGetUsersWithExistingUsers() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
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

	users, err := s.usersRepository.GetUsers(ctx, nil)
	s.NoError(err)
	s.NotEmpty(users)
}

func (s *UsersRepositoryTestSuite) TestGetUsersWithExistingUsersAndPagination() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
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

	pagination := &entities.Pagination{
		Limit:  pointers.New[uint64](1),
		Offset: pointers.New[uint64](1),
	}

	users, err := s.usersRepository.GetUsers(ctx, pagination)
	s.NoError(err)
	s.Empty(users)
}

func (s *UsersRepositoryTestSuite) TestGetUsersWithoutExistingUsers() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	users, err := s.usersRepository.GetUsers(ctx, nil)
	s.NoError(err)
	s.Empty(users)
}

func (s *UsersRepositoryTestSuite) TestUpdateUserProfileSuccess() {
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

	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	err = s.usersRepository.UpdateUserProfile(
		ctx,
		entities.UpdateUserProfileDTO{
			UserID:      testUser.ID,
			DisplayName: &testUser.DisplayName,
			Phone:       testUser.Phone,
			Telegram:    testUser.Telegram,
			Avatar:      testUser.Avatar,
		},
	)

	s.NoError(err)
}

func (s *UsersRepositoryTestSuite) TestUpdateUserProfileUserDoesNotExists() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	err := s.usersRepository.UpdateUserProfile(
		ctx,
		entities.UpdateUserProfileDTO{
			UserID:      testUser.ID,
			DisplayName: &testUser.DisplayName,
			Phone:       testUser.Phone,
			Telegram:    testUser.Telegram,
			Avatar:      testUser.Avatar,
		},
	)

	s.NoError(err)
}

func BenchmarkUsersRepository_GetUserByID(b *testing.B) {
	spanConfig := tracing.SpanConfig{}
	ctrl := gomock.NewController(b)
	logger := loggermock.NewMockLogger(ctrl)
	dbConnector, err := db.New(dsn, driver, logger)
	require.NoError(b, err)

	defer func() {
		require.NoError(b, dbConnector.Close())
	}()

	traceProvider := mocktracing.NewMockProvider(ctrl)
	traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(ctx, mocktracing.NewMockSpan()).
		AnyTimes()

	usersRepository := repositories.NewUsersRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = usersRepository.GetUserByID(
			ctx,
			userID,
		)
	}
}

func BenchmarkUsersRepository_GetUserByEmail(b *testing.B) {
	spanConfig := tracing.SpanConfig{}
	ctrl := gomock.NewController(b)
	logger := loggermock.NewMockLogger(ctrl)
	dbConnector, err := db.New(dsn, driver, logger)
	require.NoError(b, err)

	defer func() {
		require.NoError(b, dbConnector.Close())
	}()

	traceProvider := mocktracing.NewMockProvider(ctrl)
	traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(ctx, mocktracing.NewMockSpan()).
		AnyTimes()

	usersRepository := repositories.NewUsersRepository(dbConnector, logger, traceProvider, spanConfig)

	b.ResetTimer()
	for range b.N {
		_, _ = usersRepository.GetUserByEmail(
			ctx,
			email,
		)
	}
}

func BenchmarkUsersRepository_GetUsers(b *testing.B) {
	spanConfig := tracing.SpanConfig{}
	ctrl := gomock.NewController(b)
	logger := loggermock.NewMockLogger(ctrl)
	dbConnector, err := db.New(dsn, driver, logger)
	require.NoError(b, err)

	defer func() {
		require.NoError(b, dbConnector.Close())
	}()

	traceProvider := mocktracing.NewMockProvider(ctrl)
	traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(ctx, mocktracing.NewMockSpan()).
		AnyTimes()

	usersRepository := repositories.NewUsersRepository(dbConnector, logger, traceProvider, spanConfig)

	pagination := &entities.Pagination{
		Limit:  pointers.New[uint64](1),
		Offset: pointers.New[uint64](1),
	}

	b.ResetTimer()
	for range b.N {
		_, _ = usersRepository.GetUsers(ctx, pagination)
	}
}
