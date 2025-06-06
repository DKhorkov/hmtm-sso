package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	notifications "github.com/DKhorkov/hmtm-notifications/dto"
	mockcache "github.com/DKhorkov/libs/cache/mocks"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	mocknats "github.com/DKhorkov/libs/nats/mocks"
	"github.com/DKhorkov/libs/pointers"
	"github.com/DKhorkov/libs/security"
	"github.com/DKhorkov/libs/validation"

	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	mockservices "github.com/DKhorkov/hmtm-sso/mocks/services"
)

var (
	cfg              = config.New()
	validationConfig = cfg.Validation
)

func TestUseCases_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		HashCost: 10,
		JWT: security.JWTConfig{
			SecretKey: "secret",
		},
	}

	natsConfig := config.NATSConfig{
		Subjects: config.NATSSubjects{
			VerifyEmail: "verify-email",
		},
	}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		userData   entities.RegisterUserDTO
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedID  uint64
		expectedErr error
	}{
		{
			name: "success",
			userData: entities.RegisterUserDTO{
				Email:       "test@example.com",
				Password:    "Password123@",
				DisplayName: "Иван",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(uint64(1), nil).
					Times(1)

				verifyEmailDTO := notifications.VerifyEmailDTO{UserID: uint64(1)}
				content, _ := json.Marshal(verifyEmailDTO)
				natsPublisher.
					EXPECT().
					Publish("verify-email", content).
					Return(nil).
					Times(1)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "invalid email",
			userData: entities.RegisterUserDTO{
				Email:       "invalid",
				Password:    "Password123@",
				DisplayName: "Иван",
			},
			expectedID:  0,
			expectedErr: &validation.Error{},
		},
		{
			name: "invalid password",
			userData: entities.RegisterUserDTO{
				Email:       "test@example.com",
				Password:    "short",
				DisplayName: "Иван",
			},
			expectedID:  0,
			expectedErr: &validation.Error{},
		},
		{
			name: "invalid display name",
			userData: entities.RegisterUserDTO{
				Email:       "test@example.com",
				Password:    "Password123@",
				DisplayName: "",
			},
			expectedID:  0,
			expectedErr: &validation.Error{},
		},
		{
			name: "publish error",
			userData: entities.RegisterUserDTO{
				Email:       "test@example.com",
				Password:    "Password123@",
				DisplayName: "Иван",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(uint64(1), nil).
					Times(1)

				verifyEmailDTO := notifications.VerifyEmailDTO{UserID: uint64(1)}
				content, _ := json.Marshal(verifyEmailDTO)
				natsPublisher.
					EXPECT().
					Publish("verify-email", content).
					Return(errors.New("publish failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "error",
			userData: entities.RegisterUserDTO{
				Email:       "test@example.com",
				Password:    "Password123@",
				DisplayName: "Иван",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("test")).
					Times(1)
			},
			expectedErr: errors.New("test"),
		},
		{
			name: "forbidden word in display name",
			userData: entities.RegisterUserDTO{
				Email:       "test@example.com",
				Password:    "Password123@",
				DisplayName: "сука",
			},
			expectedErr: &validation.Error{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			userID, err := useCases.RegisterUser(context.Background(), tc.userData)
			require.Equal(t, tc.expectedID, userID)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:       "secret",
			Algorithm:       "HS256",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		},
		HashCost: 10,
	}
	natsConfig := config.NATSConfig{}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		userData   entities.LoginUserDTO
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name: "success",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				hashedPassword, _ := security.Hash("password123", 10)
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{
						ID:             1,
						Email:          "test@example.com",
						Password:       hashedPassword,
						EmailConfirmed: true,
					}, nil).
					Times(1)

				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				authService.
					EXPECT().
					CreateRefreshToken(
						gomock.Any(),
						uint64(1),
						gomock.Any(),
						time.Hour,
					).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "email not confirmed",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{
						ID:             1,
						Email:          "test@example.com",
						Password:       "hashed_password",
						EmailConfirmed: false,
					}, nil).
					Times(1)
			},
			expectedErr: &customerrors.EmailIsNotConfirmedError{},
		},
		{
			name: "wrong password",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				hashedPassword, _ := security.Hash("password123", 10)
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{
						ID:             1,
						Email:          "test@example.com",
						Password:       hashedPassword,
						EmailConfirmed: true,
					}, nil).
					Times(1)
			},
			expectedErr: &customerrors.WrongPasswordError{},
		},
		{
			name: "user not found",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name: "expire refresh token error",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				hashedPassword, _ := security.Hash("password123", 10)
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{
						ID:             1,
						Email:          "test@example.com",
						Password:       hashedPassword,
						EmailConfirmed: true,
					}, nil).
					Times(1)

				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(&entities.RefreshToken{}, nil).
					Times(1)

				authService.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), gomock.Any()).
					Return(errors.New("test")).
					Times(1)
			},
			expectedErr: errors.New("test"),
		},
		{
			name: "create refresh token error",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				hashedPassword, _ := security.Hash("password123", 10)
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{
						ID:             1,
						Email:          "test@example.com",
						Password:       hashedPassword,
						EmailConfirmed: true,
					}, nil).
					Times(1)

				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.RefreshTokenNotFoundError{}).
					Times(1)

				authService.
					EXPECT().
					CreateRefreshToken(
						gomock.Any(),
						uint64(1),
						gomock.Any(),
						time.Hour,
					).
					Return(uint64(0), errors.New("test")).
					Times(1)
			},
			expectedErr: errors.New("test"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			tokens, err := useCases.LoginUser(context.Background(), tc.userData)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.NotZero(t, tokens.AccessToken)
				require.NotZero(t, tokens.RefreshToken)
			}
		})
	}
}

func TestUseCases_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{}
	natsConfig := config.NATSConfig{}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		id         uint64
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedUser *entities.User
		expectedErr  error
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1}, nil).
					Times(1)
			},
			expectedUser: &entities.User{ID: 1},
			expectedErr:  nil,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedUser: nil,
			expectedErr:  errors.New("not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			user, err := useCases.GetUserByID(context.Background(), tc.id)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUseCases_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{}
	natsConfig := config.NATSConfig{}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		email      string
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedUser *entities.User
		expectedErr  error
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, Email: "test@example.com"}, nil).
					Times(1)
			},
			expectedUser: &entities.User{ID: 1, Email: "test@example.com"},
			expectedErr:  nil,
		},
		{
			name:  "error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedUser: nil,
			expectedErr:  errors.New("not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			user, err := useCases.GetUserByEmail(context.Background(), tc.email)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUseCases_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{}
	natsConfig := config.NATSConfig{}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		pagination *entities.Pagination
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedUsers []entities.User
		expectedErr   error
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return([]entities.User{{ID: 1}, {ID: 2}}, nil).
					Times(1)
			},
			expectedUsers: []entities.User{{ID: 1}, {ID: 2}},
			expectedErr:   nil,
		},
		{
			name: "error",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedUsers: nil,
			expectedErr:   errors.New("fetch failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			users, err := useCases.GetUsers(context.Background(), tc.pagination)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUsers, users)
			}
		})
	}
}

func TestUseCases_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:       "secret",
			Algorithm:       "HS256",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		},
		HashCost: 10,
	}

	natsConfig := config.NATSConfig{}

	accessToken, err := security.GenerateJWT(
		uint64(1),
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidAccessToken, err := security.GenerateJWT(
		"invalid",
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		userData   entities.RawUpdateUserProfileDTO
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name: "success",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: accessToken,
				DisplayName: pointers.New("Иван"),
				Phone:       pointers.New("89112580162"),
				Telegram:    pointers.New("@tests"),
				Avatar:      pointers.New("http://someurl"),
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1}, nil).
					Times(1)

				usersService.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.UpdateUserProfileDTO{
						UserID:      1,
						DisplayName: pointers.New("Иван"),
						Phone:       pointers.New("89112580162"),
						Telegram:    pointers.New("@tests"),
						Avatar:      pointers.New("http://someurl"),
					}).
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "invalid display name",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_token",
				DisplayName: pointers.New(""),
			},
			expectedErr: &validation.Error{},
		},
		{
			name: "invalid phone",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_token",
				DisplayName: pointers.New("Иван"),
				Phone:       pointers.New("t"),
			},
			expectedErr: &validation.Error{},
		},
		{
			name: "invalid telegram",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_token",
				DisplayName: pointers.New("Иван"),
				Telegram:    pointers.New("t"),
			},
			expectedErr: &validation.Error{},
		},
		{
			name: "invalid access token",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: "invalid_token",
				DisplayName: pointers.New("Иван"),
			},
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name: "invalid access token payload",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: invalidAccessToken,
				DisplayName: pointers.New("Иван"),
			},
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name: "get user by id error",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: accessToken,
				DisplayName: pointers.New("Иван"),
				Phone:       pointers.New("89112580162"),
				Telegram:    pointers.New("@tests"),
				Avatar:      pointers.New("http://someurl"),
			},
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name: "forbidden word in display name",
			userData: entities.RawUpdateUserProfileDTO{
				AccessToken: accessToken,
				DisplayName: pointers.New("Сука"),
				Phone:       pointers.New("89112580162"),
				Telegram:    pointers.New("@tests"),
				Avatar:      pointers.New("http://someurl"),
			},
			expectedErr: &validation.Error{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err = useCases.UpdateUserProfile(context.Background(), tc.userData)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:       "secret",
			Algorithm:       "HS256",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		},
		HashCost: 10,
	}
	natsConfig := config.NATSConfig{}

	accessToken, err := security.GenerateJWT(
		uint64(1),
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidAccessToken, err := security.GenerateJWT(
		"invalid",
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name        string
		accessToken string
		setupMocks  func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedUser *entities.User
		expectedErr  error
	}{
		{
			name:        "success",
			accessToken: accessToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1}, nil).
					Times(1)
			},
			expectedUser: &entities.User{ID: 1},
			expectedErr:  nil,
		},
		{
			name:         "invalid token",
			accessToken:  "invalid_token",
			expectedUser: nil,
			expectedErr:  &security.InvalidJWTError{},
		},
		{
			name:        "invalid token payload",
			accessToken: invalidAccessToken,
			expectedErr: &security.InvalidJWTError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			user, err := useCases.GetMe(context.Background(), tc.accessToken)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUseCases_RefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:       "secret",
			Algorithm:       "HS256",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour,
		},
		HashCost: 10,
	}
	natsConfig := config.NATSConfig{}

	accessToken, err := security.GenerateJWT(
		uint64(1),
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidAccessToken, err := security.GenerateJWT(
		"invalid",
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	refreshToken, err := security.GenerateJWT(
		accessToken,
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.RefreshTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidRefreshTokenInt, err := security.GenerateJWT(
		1,
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.RefreshTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidRefreshTokenString, err := security.GenerateJWT(
		"invalid",
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.RefreshTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidRefreshTokenWithInvalidAccessToken, err := security.GenerateJWT(
		invalidAccessToken,
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.RefreshTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	encodedRefreshToken := security.RawEncode([]byte(refreshToken))

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name         string
		refreshToken string
		setupMocks   func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:         "success",
			refreshToken: encodedRefreshToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(&entities.RefreshToken{Value: refreshToken}, nil).
					Times(1)

				authService.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				authService.
					EXPECT().
					CreateRefreshToken(
						gomock.Any(),
						uint64(1),
						gomock.Any(),
						time.Hour,
					).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid_token",
			expectedErr:  &security.InvalidJWTError{},
		},
		{
			name:         "invalid refresh token payload",
			refreshToken: security.RawEncode([]byte("invalid")),
			expectedErr:  &security.InvalidJWTError{},
		},
		{
			name:         "invalid refresh token after encoding",
			refreshToken: security.RawEncode([]byte(invalidRefreshTokenInt)),
			expectedErr:  &security.InvalidJWTError{},
		},
		{
			name:         "invalid access token payload",
			refreshToken: security.RawEncode([]byte(invalidRefreshTokenString)),
			expectedErr:  &security.InvalidJWTError{},
		},
		{
			name:         "invalid access token payload",
			refreshToken: security.RawEncode([]byte(invalidRefreshTokenWithInvalidAccessToken)),
			expectedErr:  &security.InvalidJWTError{},
		},
		{
			name:         "get db refresh token by id error",
			refreshToken: encodedRefreshToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.RefreshTokenNotFoundError{}).
					Times(1)
			},
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name:         "expire db refresh token error",
			refreshToken: encodedRefreshToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(&entities.RefreshToken{Value: refreshToken}, nil).
					Times(1)

				authService.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), gomock.Any()).
					Return(errors.New("test")).
					Times(1)
			},
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name:         "create db refresh token error",
			refreshToken: encodedRefreshToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(&entities.RefreshToken{Value: refreshToken}, nil).
					Times(1)

				authService.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				authService.
					EXPECT().
					CreateRefreshToken(
						gomock.Any(),
						uint64(1),
						gomock.Any(),
						time.Hour,
					).
					Return(uint64(0), errors.New("test")).
					Times(1)
			},
			expectedErr: errors.New("test"),
		},
		{
			name:         "access token does not belong to refresh token",
			refreshToken: encodedRefreshToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(&entities.RefreshToken{Value: invalidRefreshTokenString}, nil).
					Times(1)
			},
			expectedErr: &customerrors.AccessTokenDoesNotBelongToRefreshTokenError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			tokens, err := useCases.RefreshTokens(context.Background(), tc.refreshToken)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.NotZero(t, tokens.AccessToken)
				require.NotZero(t, tokens.RefreshToken)
			}
		})
	}
}

func TestUseCases_LogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:      "secret",
			Algorithm:      "HS256",
			AccessTokenTTL: time.Hour,
		},
		HashCost: 10,
	}
	natsConfig := config.NATSConfig{}

	accessToken, err := security.GenerateJWT(
		uint64(1),
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidAccessToken, err := security.GenerateJWT(
		"invalid",
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name        string
		accessToken string
		setupMocks  func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:        "success",
			accessToken: accessToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(&entities.RefreshToken{Value: "refresh_token"}, nil).
					Times(1)

				authService.
					EXPECT().
					ExpireRefreshToken(gomock.Any(), "refresh_token").
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:        "no refresh token",
			accessToken: accessToken,
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				authService.
					EXPECT().
					GetRefreshTokenByUserID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:        "invalid token",
			accessToken: "invalid_token",
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name:        "invalid token payload",
			accessToken: invalidAccessToken,
			expectedErr: &security.InvalidJWTError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err = useCases.LogoutUser(context.Background(), tc.accessToken)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_VerifyUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{}
	natsConfig := config.NATSConfig{}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name             string
		verifyEmailToken string
		setupMocks       func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:             "success",
			verifyEmailToken: security.RawEncode([]byte("1")),
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, EmailConfirmed: false}, nil).
					Times(1)

				authService.
					EXPECT().
					VerifyUserEmail(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:             "email already confirmed",
			verifyEmailToken: security.RawEncode([]byte("1")),
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, EmailConfirmed: true}, nil).
					Times(1)
			},
			expectedErr: &customerrors.EmailAlreadyConfirmedError{},
		},
		{
			name:             "invalid token payload",
			verifyEmailToken: security.RawEncode([]byte(nil)),
			expectedErr:      &strconv.NumError{},
		},
		{
			name:             "user not found",
			verifyEmailToken: security.RawEncode([]byte("1")),
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err := useCases.VerifyUserEmail(context.Background(), tc.verifyEmailToken)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_ForgetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:      "secret",
			Algorithm:      "HS256",
			AccessTokenTTL: time.Hour,
		},
		HashCost: 10,
	}
	natsConfig := config.NATSConfig{}

	forgetPasswordToken := security.RawEncode([]byte("1"))

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name        string
		token       string
		newPassword string
		setupMocks  func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:        "success",
			token:       forgetPasswordToken,
			newPassword: "Password123@",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				oldHashedPassword, _ := security.Hash("oldpassword123", 10)
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, Password: oldHashedPassword}, nil).
					Times(1)

				authService.
					EXPECT().
					ForgetPassword(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:        "invalid password",
			token:       forgetPasswordToken,
			newPassword: "short",
			expectedErr: &validation.Error{},
		},
		{
			name:        "invalid token conversion to int",
			token:       security.RawEncode([]byte("s")),
			newPassword: "Password123@",
			expectedErr: &strconv.NumError{},
		},
		{
			name:        "get user by id error",
			token:       forgetPasswordToken,
			newPassword: "Password123@",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:        "new password is equal to old password",
			token:       forgetPasswordToken,
			newPassword: "Password123@",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				oldHashedPassword, _ := security.Hash("Password123@", 10)
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, Password: oldHashedPassword}, nil).
					Times(1)
			},
			expectedErr: &validation.Error{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err := useCases.ForgetPassword(context.Background(), tc.token, tc.newPassword)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{
		JWT: security.JWTConfig{
			SecretKey:      "secret",
			Algorithm:      "HS256",
			AccessTokenTTL: time.Hour,
		},
		HashCost: 10,
	}
	natsConfig := config.NATSConfig{}

	accessToken, err := security.GenerateJWT(
		uint64(1),
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	invalidAccessToken, err := security.GenerateJWT(
		"invalid",
		securityConfig.JWT.SecretKey,
		securityConfig.JWT.AccessTokenTTL,
		securityConfig.JWT.Algorithm,
	)
	require.NoError(t, err)

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name        string
		accessToken string
		oldPassword string
		newPassword string
		setupMocks  func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:        "success",
			accessToken: accessToken,
			oldPassword: "oldpassword123",
			newPassword: "Password123@",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				hashedPassword, _ := security.Hash("oldpassword123", 10)
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, Password: hashedPassword}, nil).
					Times(1)

				authService.
					EXPECT().
					ChangePassword(gomock.Any(), uint64(1), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:        "same password",
			accessToken: accessToken,
			oldPassword: "Password123@",
			newPassword: "Password123@",
			expectedErr: &validation.Error{},
		},
		{
			name:        "wrong old password",
			accessToken: accessToken,
			oldPassword: "wrongpassword",
			newPassword: "Password123@",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				hashedPassword, _ := security.Hash("oldpassword123", 10)
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, Password: hashedPassword}, nil).
					Times(1)
			},
			expectedErr: &customerrors.WrongPasswordError{},
		},
		{
			name:        "invalid password",
			accessToken: accessToken,
			oldPassword: "password123",
			newPassword: "safa",
			expectedErr: &validation.Error{},
		},
		{
			name:        "invalid accessToken",
			accessToken: "invalid",
			oldPassword: "oldpassword123",
			newPassword: "Password123@",
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name:        "failed to parse accessToken",
			accessToken: invalidAccessToken,
			oldPassword: "oldpassword123",
			newPassword: "Password123@",
			expectedErr: &security.InvalidJWTError{},
		},
		{
			name:        "get user by id error",
			accessToken: accessToken,
			oldPassword: "oldpassword123",
			newPassword: "Password123@",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				usersService.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err = useCases.ChangePassword(context.Background(), tc.accessToken, tc.oldPassword, tc.newPassword)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_SendVerifyEmailMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{}
	natsConfig := config.NATSConfig{
		Subjects: config.NATSSubjects{
			VerifyEmail: "verify-email",
		},
	}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		email      string
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: false}, nil).
					Times(1)

				verifyEmailDTO := notifications.VerifyEmailDTO{UserID: uint64(1)}
				content, _ := json.Marshal(verifyEmailDTO)
				natsPublisher.
					EXPECT().
					Publish("verify-email", content).
					Return(nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Incr(gomock.Any(), gomock.Any()).
					Return(int64(1), nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "email already confirmed",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: true}, nil).
					Times(1)
			},
			expectedErr: &customerrors.EmailAlreadyConfirmedError{},
		},
		{
			name:  "publish error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: false}, nil).
					Times(1)

				verifyEmailDTO := notifications.VerifyEmailDTO{UserID: uint64(1)}
				content, _ := json.Marshal(verifyEmailDTO)
				natsPublisher.
					EXPECT().
					Publish("verify-email", content).
					Return(errors.New("publish failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr: errors.New("publish failed"),
		},
		{
			name:  "user not found",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:  "cache incr error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: false}, nil).
					Times(1)

				forgetPasswordDTO := notifications.ForgetPasswordDTO{UserID: uint64(1)}
				content, _ := json.Marshal(forgetPasswordDTO)
				natsPublisher.
					EXPECT().
					Publish("verify-email", content).
					Return(nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Incr(gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("cache incr error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:  "get cache error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("", errors.New("")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:  "invalid cache error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("invalid", nil).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:  "limit exceeded error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("4", nil).
					Times(1)

			},
			expectedErr: &customerrors.LimitExceededError{},
		},
		{
			name:  "cache set error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com")).
					Return("0", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: false}, nil).
					Times(1)

				forgetPasswordDTO := notifications.ForgetPasswordDTO{UserID: uint64(1)}
				content, _ := json.Marshal(forgetPasswordDTO)
				natsPublisher.
					EXPECT().
					Publish("verify-email", content).
					Return(nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Set(
						gomock.Any(),
						fmt.Sprintf("%s-%s", verifyEmailCachePrefix, "test@example.com"),
						1,
						verifyEmailTTL,
					).
					Return(errors.New("cache incr error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err := useCases.SendVerifyEmailMessage(context.Background(), tc.email)
			if tc.expectedErr != nil {
				require.Error(t, err)
				if tc.expectedErr != nil {
					require.Error(t, err)
					require.IsType(t, tc.expectedErr, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestUseCases_SendForgetPasswordMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	authService := mockservices.NewMockAuthService(ctrl)
	usersService := mockservices.NewMockUsersService(ctrl)
	natsPublisher := mocknats.NewMockPublisher(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	cacheProvider := mockcache.NewMockProvider(ctrl)

	securityConfig := security.Config{}
	natsConfig := config.NATSConfig{
		Subjects: config.NATSSubjects{
			ForgetPassword: "forget-password",
		},
	}

	useCases := New(
		authService,
		usersService,
		securityConfig,
		validationConfig,
		natsPublisher,
		natsConfig,
		logger,
		cacheProvider,
	)

	testCases := []struct {
		name       string
		email      string
		setupMocks func(
			authService *mockservices.MockAuthService,
			usersService *mockservices.MockUsersService,
			natsPublisher *mocknats.MockPublisher,
			logger *mocklogging.MockLogger,
			cacheProvider *mockcache.MockProvider,
		)
		expectedErr error
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: true}, nil).
					Times(1)

				forgetPasswordDTO := notifications.ForgetPasswordDTO{UserID: uint64(1)}
				content, _ := json.Marshal(forgetPasswordDTO)
				natsPublisher.
					EXPECT().
					Publish("forget-password", content).
					Return(nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Incr(gomock.Any(), gomock.Any()).
					Return(int64(1), nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "email not confirmed",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: false}, nil).
					Times(1)
			},
			expectedErr: &customerrors.EmailIsNotConfirmedError{},
		},
		{
			name:  "user not found",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:  "publish error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: true}, nil).
					Times(1)

				forgetPasswordDTO := notifications.ForgetPasswordDTO{UserID: uint64(1)}
				content, _ := json.Marshal(forgetPasswordDTO)
				natsPublisher.
					EXPECT().
					Publish("forget-password", content).
					Return(errors.New("publish failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr: errors.New("publish failed"),
		},
		{
			name:  "cache incr error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("1", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: true}, nil).
					Times(1)

				forgetPasswordDTO := notifications.ForgetPasswordDTO{UserID: uint64(1)}
				content, _ := json.Marshal(forgetPasswordDTO)
				natsPublisher.
					EXPECT().
					Publish("forget-password", content).
					Return(nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Incr(gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("cache incr error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:  "get cache error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("", errors.New("")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:  "invalid cache error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("invalid", nil).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(nil, &customerrors.UserNotFoundError{}).
					Times(1)
			},
			expectedErr: &customerrors.UserNotFoundError{},
		},
		{
			name:  "limit exceeded error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("4", nil).
					Times(1)

			},
			expectedErr: &customerrors.LimitExceededError{},
		},
		{
			name:  "cache set error",
			email: "test@example.com",
			setupMocks: func(
				authService *mockservices.MockAuthService,
				usersService *mockservices.MockUsersService,
				natsPublisher *mocknats.MockPublisher,
				logger *mocklogging.MockLogger,
				cacheProvider *mockcache.MockProvider,
			) {
				cacheProvider.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Get(gomock.Any(), fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com")).
					Return("0", nil).
					Times(1)

				usersService.
					EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&entities.User{ID: 1, EmailConfirmed: true}, nil).
					Times(1)

				forgetPasswordDTO := notifications.ForgetPasswordDTO{UserID: uint64(1)}
				content, _ := json.Marshal(forgetPasswordDTO)
				natsPublisher.
					EXPECT().
					Publish("forget-password", content).
					Return(nil).
					Times(1)

				cacheProvider.
					EXPECT().
					Set(
						gomock.Any(),
						fmt.Sprintf("%s-%s", forgetPasswordCachePrefix, "test@example.com"),
						1,
						forgetPasswordTTL,
					).
					Return(errors.New("cache incr error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					authService,
					usersService,
					natsPublisher,
					logger,
					cacheProvider,
				)
			}

			err := useCases.SendForgetPasswordMessage(context.Background(), tc.email)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
