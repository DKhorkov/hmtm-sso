package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	mockrepositories "github.com/DKhorkov/hmtm-sso/mocks/repositories"
)

func TestUsersService_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewUsersService(usersRepository, logger)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		setupMocks    func(usersRepository *mockrepositories.MockUsersRepository)
		expectedUsers []entities.User
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository) {
				users := []entities.User{
					{ID: 1, Email: "user1@example.com"},
					{ID: 2, Email: "user2@example.com"},
				}
				usersRepository.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(users, nil).
					Times(1)
			},
			expectedUsers: []entities.User{
				{ID: 1, Email: "user1@example.com"},
				{ID: 2, Email: "user2@example.com"},
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "repo error",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository) {
				usersRepository.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedUsers: nil,
			expectedErr:   errors.New("database error"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usersRepository)
			}

			users, err := service.GetUsers(context.Background(), tc.pagination)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, users)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUsers, users)
			}
		})
	}
}

func TestUsersService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewUsersService(usersRepository, logger)

	testCases := []struct {
		name          string
		userID        uint64
		setupMocks    func(usersRepository *mockrepositories.MockUsersRepository, logger *mocklogging.MockLogger)
		expectedUser  *entities.User
		expectedErr   error
		errorExpected bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository, logger *mocklogging.MockLogger) {
				user := &entities.User{ID: 1, Email: "user1@example.com"}
				usersRepository.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(user, nil).
					Times(1)
			},
			expectedUser:  &entities.User{ID: 1, Email: "user1@example.com"},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:   "not found",
			userID: 1,
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository, logger *mocklogging.MockLogger) {
				usersRepository.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("user not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUser:  nil,
			expectedErr:   &customerrors.UserNotFoundError{},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usersRepository, logger)
			}

			user, err := service.GetUserByID(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUsersService_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewUsersService(usersRepository, logger)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(usersRepository *mockrepositories.MockUsersRepository, logger *mocklogging.MockLogger)
		expectedUser  *entities.User
		expectedErr   error
		errorExpected bool
	}{
		{
			name:  "success",
			email: "user1@example.com",
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository, logger *mocklogging.MockLogger) {
				user := &entities.User{ID: 1, Email: "user1@example.com"}
				usersRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "user1@example.com").
					Return(user, nil).
					Times(1)
			},
			expectedUser:  &entities.User{ID: 1, Email: "user1@example.com"},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:  "not found",
			email: "user1@example.com",
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository, logger *mocklogging.MockLogger) {
				usersRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "user1@example.com").
					Return(nil, errors.New("user not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUser:  nil,
			expectedErr:   &customerrors.UserNotFoundError{},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usersRepository, logger)
			}

			user, err := service.GetUserByEmail(context.Background(), tc.email)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUsersService_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewUsersService(usersRepository, logger)

	testCases := []struct {
		name            string
		userProfileData entities.UpdateUserProfileDTO
		setupMocks      func(usersRepository *mockrepositories.MockUsersRepository)
		expectedErr     error
		errorExpected   bool
	}{
		{
			name: "success",
			userProfileData: entities.UpdateUserProfileDTO{
				UserID:      1,
				DisplayName: pointers.New("name"),
				Phone:       pointers.New("89112580162"),
				Telegram:    pointers.New("@test"),
				Avatar:      pointers.New("http://someurl"),
			},
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository) {
				usersRepository.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						entities.UpdateUserProfileDTO{
							UserID:      1,
							DisplayName: pointers.New("name"),
							Phone:       pointers.New("89112580162"),
							Telegram:    pointers.New("@test"),
							Avatar:      pointers.New("http://someurl"),
						},
					).
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "repo error",
			userProfileData: entities.UpdateUserProfileDTO{
				UserID:      1,
				DisplayName: pointers.New("name"),
				Phone:       pointers.New("89112580162"),
				Telegram:    pointers.New("@test"),
				Avatar:      pointers.New("http://someurl"),
			},
			setupMocks: func(usersRepository *mockrepositories.MockUsersRepository) {
				usersRepository.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						entities.UpdateUserProfileDTO{
							UserID:      1,
							DisplayName: pointers.New("name"),
							Phone:       pointers.New("89112580162"),
							Telegram:    pointers.New("@test"),
							Avatar:      pointers.New("http://someurl"),
						},
					).
					Return(errors.New("update failed")).
					Times(1)
			},
			expectedErr:   errors.New("update failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usersRepository)
			}

			err := service.UpdateUserProfile(context.Background(), tc.userProfileData)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
