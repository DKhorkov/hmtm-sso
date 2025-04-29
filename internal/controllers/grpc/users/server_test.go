package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	customgrpc "github.com/DKhorkov/libs/grpc"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	mockusecases "github.com/DKhorkov/hmtm-sso/mocks/usecases"
	"github.com/DKhorkov/libs/security"
)

func TestServerAPI_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.UpdateUserProfileIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in: &sso.UpdateUserProfileIn{
				AccessToken: "valid-token",
				DisplayName: pointers.New("John Doe"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@johndoe"),
				Avatar:      pointers.New("avatar.jpg"),
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.RawUpdateUserProfileDTO{
						AccessToken: "valid-token",
						DisplayName: pointers.New("John Doe"),
						Phone:       pointers.New("1234567890"),
						Telegram:    pointers.New("@johndoe"),
						Avatar:      pointers.New("avatar.jpg"),
					}).
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "invalid phone",
			in: &sso.UpdateUserProfileIn{
				AccessToken: "valid-token",
				Phone:       pointers.New("invalid"),
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.RawUpdateUserProfileDTO{
						AccessToken: "valid-token",
						Phone:       pointers.New("invalid"),
					}).
					Return(&customerrors.InvalidPhoneError{Message: "phone format invalid"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: "phone format invalid"},
			errorExpected: true,
		},
		{
			name: "invalid JWT",
			in:   &sso.UpdateUserProfileIn{AccessToken: "invalid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.RawUpdateUserProfileDTO{
						AccessToken: "invalid-token",
					}).
					Return(&security.InvalidJWTError{Message: "token invalid"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Unauthenticated, Message: "token invalid"},
			errorExpected: true,
		},
		{
			name: "user not found",
			in:   &sso.UpdateUserProfileIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.RawUpdateUserProfileDTO{
						AccessToken: "valid-token",
					}).
					Return(&customerrors.UserNotFoundError{Message: "user not found"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.NotFound, Message: "user not found"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.UpdateUserProfileIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.RawUpdateUserProfileDTO{
						AccessToken: "valid-token",
					}).
					Return(errors.New("internal error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Internal, Message: "internal error"},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			resp, err := api.UpdateUserProfile(context.Background(), tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.IsType(t, &emptypb.Empty{}, resp)
			}
		})
	}
}

func TestServerAPI_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.GetUserByEmailIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.GetUserOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.GetUserByEmailIn{Email: "john@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				user := &entities.User{
					ID:          1,
					DisplayName: "John Doe",
					Email:       "john@example.com",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				}
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), "john@example.com").
					Return(user, nil).
					Times(1)
			},
			expectedOut: &sso.GetUserOut{
				ID:          1,
				DisplayName: "John Doe",
				Email:       "john@example.com",
				CreatedAt:   timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:   timestamppb.New(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user not found",
			in:   &sso.GetUserByEmailIn{Email: "john@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), "john@example.com").
					Return(nil, &customerrors.UserNotFoundError{Message: "user not found"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.NotFound, Message: "user not found"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.GetUserByEmailIn{Email: "john@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), "john@example.com").
					Return(nil, errors.New("internal error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Internal, Message: "internal error"},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			resp, err := api.GetUserByEmail(context.Background(), tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOut, resp)
			}
		})
	}
}

func TestServerAPI_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.GetUserIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.GetUserOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.GetUserIn{ID: 1},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				user := &entities.User{
					ID:          1,
					DisplayName: "John Doe",
					Email:       "john@example.com",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				}
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(user, nil).
					Times(1)
			},
			expectedOut: &sso.GetUserOut{
				ID:          1,
				DisplayName: "John Doe",
				Email:       "john@example.com",
				CreatedAt:   timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:   timestamppb.New(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user not found",
			in:   &sso.GetUserIn{ID: 1},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.UserNotFoundError{Message: "user not found"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.NotFound, Message: "user not found"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.GetUserIn{ID: 1},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("internal error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Internal, Message: "internal error"},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			resp, err := api.GetUser(context.Background(), tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOut, resp)
			}
		})
	}
}

func TestServerAPI_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.GetUsersOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				users := []entities.User{
					{
						ID:          1,
						DisplayName: "John Doe",
						Email:       "john@example.com",
						CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          2,
						DisplayName: "Jane Doe",
						Email:       "jane@example.com",
						CreatedAt:   time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC),
					},
				}
				useCases.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return(users, nil).
					Times(1)
			},
			expectedOut: &sso.GetUsersOut{
				Users: []*sso.GetUserOut{
					{
						ID:          1,
						DisplayName: "John Doe",
						Email:       "john@example.com",
						CreatedAt:   timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
						UpdatedAt:   timestamppb.New(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:          2,
						DisplayName: "Jane Doe",
						Email:       "jane@example.com",
						CreatedAt:   timestamppb.New(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
						UpdatedAt:   timestamppb.New(time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "internal error",
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return(nil, errors.New("internal error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Internal, Message: "internal error"},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			resp, err := api.GetUsers(context.Background(), &emptypb.Empty{})
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOut, resp)
			}
		})
	}
}

func TestServerAPI_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.GetMeIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.GetUserOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.GetMeIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				user := &entities.User{
					ID:          1,
					DisplayName: "John Doe",
					Email:       "john@example.com",
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				}
				useCases.
					EXPECT().
					GetMe(gomock.Any(), "valid-token").
					Return(user, nil).
					Times(1)
			},
			expectedOut: &sso.GetUserOut{
				ID:          1,
				DisplayName: "John Doe",
				Email:       "john@example.com",
				CreatedAt:   timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:   timestamppb.New(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "invalid JWT",
			in:   &sso.GetMeIn{AccessToken: "invalid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetMe(gomock.Any(), "invalid-token").
					Return(nil, &security.InvalidJWTError{Message: "token invalid"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Unauthenticated, Message: "token invalid"},
			errorExpected: true,
		},
		{
			name: "user not found",
			in:   &sso.GetMeIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetMe(gomock.Any(), "valid-token").
					Return(nil, &customerrors.UserNotFoundError{Message: "user not found"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.NotFound, Message: "user not found"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.GetMeIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					GetMe(gomock.Any(), "valid-token").
					Return(nil, errors.New("internal error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Internal, Message: "internal error"},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			resp, err := api.GetMe(context.Background(), tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.expectedErr, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOut, resp)
			}
		})
	}
}
