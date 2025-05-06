package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocklogging "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	mockrepositories "github.com/DKhorkov/hmtm-sso/mocks/repositories"
)

func TestAuthService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		userData      entities.RegisterUserDTO
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository, usersRepository *mockrepositories.MockUsersRepository)
		expectedID    uint64
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			userData: entities.RegisterUserDTO{
				Email: "test@example.com",
			},
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository, usersRepository *mockrepositories.MockUsersRepository) {
				usersRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, nil).
					Times(1)

				authRepository.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{Email: "test@example.com"}).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedID:    1,
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user already exists",
			userData: entities.RegisterUserDTO{
				Email: "existing@example.com",
			},
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository, usersRepository *mockrepositories.MockUsersRepository) {
				usersRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "existing@example.com").
					Return(&entities.User{ID: 1, Email: "existing@example.com"}, nil).
					Times(1)
			},
			expectedID:    0,
			expectedErr:   &customerrors.UserAlreadyExistsError{},
			errorExpected: true,
		},
		{
			name: "auth repo error",
			userData: entities.RegisterUserDTO{
				Email: "test@example.com",
			},
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository, usersRepository *mockrepositories.MockUsersRepository) {
				usersRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, nil).
					Times(1)

				authRepository.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{Email: "test@example.com"}).
					Return(uint64(0), errors.New("registration failed")).
					Times(1)
			},
			expectedID:    0,
			expectedErr:   errors.New("registration failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository, usersRepository)
			}

			id, err := service.RegisterUser(context.Background(), tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Equal(t, tc.expectedID, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestAuthService_CreateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		userID        uint64
		refreshToken  string
		ttl           time.Duration
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository)
		expectedID    uint64
		expectedErr   error
		errorExpected bool
	}{
		{
			name:         "success",
			userID:       1,
			refreshToken: "token123",
			ttl:          time.Hour,
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					CreateRefreshToken(gomock.Any(), uint64(1), "token123", time.Hour).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedID:    1,
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:         "repo error",
			userID:       1,
			refreshToken: "token123",
			ttl:          time.Hour,
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					CreateRefreshToken(gomock.Any(), uint64(1), "token123", time.Hour).
					Return(uint64(0), errors.New("token creation failed")).
					Times(1)
			},
			expectedID:    0,
			expectedErr:   errors.New("token creation failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository)
			}

			id, err := service.CreateRefreshToken(context.Background(), tc.userID, tc.refreshToken, tc.ttl)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Equal(t, tc.expectedID, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestAuthService_GetRefreshTokenByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		userID        uint64
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository)
		expectedToken *entities.RefreshToken
		expectedErr   error
		errorExpected bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				token := &entities.RefreshToken{ID: 1, UserID: 1, Value: "token123"}
				authRepository.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(token, nil).
					Times(1)
			},
			expectedToken: &entities.RefreshToken{ID: 1, UserID: 1, Value: "token123"},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:   "repo error",
			userID: 1,
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("token not found")).
					Times(1)
			},
			expectedToken: nil,
			expectedErr:   errors.New("token not found"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository)
			}

			token, err := service.GetRefreshTokenByUserID(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
				require.Nil(t, token)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToken, token)
			}
		})
	}
}

func TestAuthService_ExpireRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		refreshToken  string
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository)
		expectedErr   error
		errorExpected bool
	}{
		{
			name:         "success",
			refreshToken: "token123",
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), "token123").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:         "repo error",
			refreshToken: "token123",
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), "token123").
					Return(errors.New("expire failed")).
					Times(1)
			},
			expectedErr:   errors.New("expire failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository)
			}

			err := service.ExpireRefreshToken(context.Background(), tc.refreshToken)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_VerifyUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		userID        uint64
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository)
		expectedErr   error
		errorExpected bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					VerifyUserEmail(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:   "repo error",
			userID: 1,
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					VerifyUserEmail(gomock.Any(), uint64(1)).
					Return(errors.New("verification failed")).
					Times(1)
			},
			expectedErr:   errors.New("verification failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository)
			}

			err := service.VerifyUserEmail(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_ForgetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		userID        uint64
		newPassword   string
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository)
		expectedErr   error
		errorExpected bool
	}{
		{
			name:        "success",
			userID:      1,
			newPassword: "newpass123",
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					ForgetPassword(gomock.Any(), uint64(1), "newpass123").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:        "repo error",
			userID:      1,
			newPassword: "newpass123",
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					ForgetPassword(gomock.Any(), uint64(1), "newpass123").
					Return(errors.New("reset failed")).
					Times(1)
			},
			expectedErr:   errors.New("reset failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository)
			}

			err := service.ForgetPassword(context.Background(), tc.userID, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	authRepository := mockrepositories.NewMockAuthRepository(ctrl)
	usersRepository := mockrepositories.NewMockUsersRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewAuthService(authRepository, usersRepository, logger)

	testCases := []struct {
		name          string
		userID        uint64
		newPassword   string
		setupMocks    func(authRepository *mockrepositories.MockAuthRepository)
		expectedErr   error
		errorExpected bool
	}{
		{
			name:        "success",
			userID:      1,
			newPassword: "newpass123",
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					ChangePassword(gomock.Any(), uint64(1), "newpass123").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name:        "repo error",
			userID:      1,
			newPassword: "newpass123",
			setupMocks: func(authRepository *mockrepositories.MockAuthRepository) {
				authRepository.
					EXPECT().
					ChangePassword(gomock.Any(), uint64(1), "newpass123").
					Return(errors.New("change failed")).
					Times(1)
			},
			expectedErr:   errors.New("change failed"),
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(authRepository)
			}

			err := service.ChangePassword(context.Background(), tc.userID, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
