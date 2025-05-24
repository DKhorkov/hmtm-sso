package auth

import (
	"context"
	"errors"
	"github.com/DKhorkov/libs/validation"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	customgrpc "github.com/DKhorkov/libs/grpc"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	mockusecases "github.com/DKhorkov/hmtm-sso/mocks/usecases"
	"github.com/DKhorkov/libs/security"
)

func TestServerAPI_SendForgetPasswordMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.SendForgetPasswordMessageIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.SendForgetPasswordMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), "test@example.com").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user not found",
			in:   &sso.SendForgetPasswordMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), "test@example.com").
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
			name: "email not confirmed",
			in:   &sso.SendForgetPasswordMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), "test@example.com").
					Return(&customerrors.EmailIsNotConfirmedError{Message: "email not confirmed"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: "email not confirmed"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.SendForgetPasswordMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), "test@example.com").
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

			resp, err := api.SendForgetPasswordMessage(context.Background(), tc.in)
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

func TestServerAPI_SendVerifyEmailMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.SendVerifyEmailMessageIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.SendVerifyEmailMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), "test@example.com").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user not found",
			in:   &sso.SendVerifyEmailMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), "test@example.com").
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
			name: "email already confirmed",
			in:   &sso.SendVerifyEmailMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), "test@example.com").
					Return(&customerrors.EmailAlreadyConfirmedError{Message: "email already confirmed"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: "email already confirmed"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.SendVerifyEmailMessageIn{Email: "test@example.com"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), "test@example.com").
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

			resp, err := api.SendVerifyEmailMessage(context.Background(), tc.in)
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

func TestServerAPI_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.ChangePasswordIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in: &sso.ChangePasswordIn{
				AccessToken: "valid-token",
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					ChangePassword(gomock.Any(), "valid-token", "oldpass", "newpass").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "internal error",
			in: &sso.ChangePasswordIn{
				AccessToken: "valid-token",
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					ChangePassword(gomock.Any(), "valid-token", "oldpass", "newpass").
					Return(errors.New("internal error")).
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

			resp, err := api.ChangePassword(context.Background(), tc.in)
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

func TestServerAPI_ForgetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.ForgetPasswordIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.ForgetPasswordIn{ForgetPasswordToken: "valid-token", NewPassword: "newpass"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					ForgetPassword(gomock.Any(), "valid-token", "newpass").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user not found",
			in:   &sso.ForgetPasswordIn{ForgetPasswordToken: "valid-token", NewPassword: "newpass"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					ForgetPassword(gomock.Any(), "valid-token", "newpass").
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
			name: "invalid password",
			in:   &sso.ForgetPasswordIn{ForgetPasswordToken: "valid-token", NewPassword: "weak"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					ForgetPassword(gomock.Any(), "valid-token", "weak").
					Return(&validation.Error{Message: "password too weak"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: "password too weak"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in:   &sso.ForgetPasswordIn{ForgetPasswordToken: "valid-token", NewPassword: "newpass"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					ForgetPassword(gomock.Any(), "valid-token", "newpass").
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

			resp, err := api.ForgetPassword(context.Background(), tc.in)
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

func TestServerAPI_VerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.VerifyEmailIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.VerifyEmailIn{VerifyEmailToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					VerifyUserEmail(gomock.Any(), "valid-token").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "email already confirmed",
			in:   &sso.VerifyEmailIn{VerifyEmailToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					VerifyUserEmail(gomock.Any(), "valid-token").
					Return(&customerrors.EmailAlreadyConfirmedError{Message: "email already confirmed"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: "email already confirmed"},
			errorExpected: true,
		},
		{
			name: "user not found",
			in:   &sso.VerifyEmailIn{VerifyEmailToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					VerifyUserEmail(gomock.Any(), "valid-token").
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
			in:   &sso.VerifyEmailIn{VerifyEmailToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					VerifyUserEmail(gomock.Any(), "valid-token").
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

			resp, err := api.VerifyEmail(context.Background(), tc.in)
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

func TestServerAPI_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.LogoutIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.LogoutIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					LogoutUser(gomock.Any(), "valid-token").
					Return(nil).
					Times(1)
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "internal error",
			in:   &sso.LogoutIn{AccessToken: "valid-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					LogoutUser(gomock.Any(), "valid-token").
					Return(errors.New("internal error")).
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

			resp, err := api.Logout(context.Background(), tc.in)
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

func TestServerAPI_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.RegisterIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.RegisterOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in: &sso.RegisterIn{
				DisplayName: "John Doe",
				Email:       "john@example.com",
				Password:    "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{
						DisplayName: "John Doe",
						Email:       "john@example.com",
						Password:    "password123",
					}).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedOut:   &sso.RegisterOut{UserID: 1},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user already exists",
			in: &sso.RegisterIn{
				DisplayName: "John Doe",
				Email:       "john@example.com",
				Password:    "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{
						DisplayName: "John Doe",
						Email:       "john@example.com",
						Password:    "password123",
					}).
					Return(uint64(0), &customerrors.UserAlreadyExistsError{Message: "user already exists"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.AlreadyExists, Message: "user already exists"},
			errorExpected: true,
		},
		{
			name: "invalid email",
			in: &sso.RegisterIn{
				DisplayName: "John Doe",
				Email:       "invalid",
				Password:    "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{
						DisplayName: "John Doe",
						Email:       "invalid",
						Password:    "password123",
					}).
					Return(uint64(0), &validation.Error{Message: "invalid email"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: "invalid email"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in: &sso.RegisterIn{
				DisplayName: "John Doe",
				Email:       "john@example.com",
				Password:    "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{
						DisplayName: "John Doe",
						Email:       "john@example.com",
						Password:    "password123",
					}).
					Return(uint64(0), errors.New("internal error")).
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

			resp, err := api.Register(context.Background(), tc.in)
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

func TestServerAPI_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.LoginIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.LoginOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in: &sso.LoginIn{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				tokens := &entities.TokensDTO{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}
				useCases.
					EXPECT().
					LoginUser(gomock.Any(), entities.LoginUserDTO{
						Email:    "john@example.com",
						Password: "password123",
					}).
					Return(tokens, nil).
					Times(1)
			},
			expectedOut: &sso.LoginOut{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "user not found",
			in: &sso.LoginIn{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					LoginUser(gomock.Any(), entities.LoginUserDTO{
						Email:    "john@example.com",
						Password: "password123",
					}).
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
			name: "wrong password",
			in: &sso.LoginIn{
				Email:    "john@example.com",
				Password: "wrongpass",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					LoginUser(gomock.Any(), entities.LoginUserDTO{
						Email:    "john@example.com",
						Password: "wrongpass",
					}).
					Return(nil, &customerrors.WrongPasswordError{Message: "wrong password"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Unauthenticated, Message: "wrong password"},
			errorExpected: true,
		},
		{
			name: "internal error",
			in: &sso.LoginIn{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					LoginUser(gomock.Any(), entities.LoginUserDTO{
						Email:    "john@example.com",
						Password: "password123",
					}).
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

			resp, err := api.Login(context.Background(), tc.in)
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

func TestServerAPI_RefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	api := &ServerAPI{
		useCases: useCases,
		logger:   logger,
	}

	testCases := []struct {
		name          string
		in            *sso.RefreshTokensIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger)
		expectedOut   *sso.LoginOut
		expectedErr   error
		errorExpected bool
	}{
		{
			name: "success",
			in:   &sso.RefreshTokensIn{RefreshToken: "valid-refresh-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				tokens := &entities.TokensDTO{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), "valid-refresh-token").
					Return(tokens, nil).
					Times(1)
			},
			expectedOut: &sso.LoginOut{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			},
			expectedErr:   nil,
			errorExpected: false,
		},
		{
			name: "invalid JWT",
			in:   &sso.RefreshTokensIn{RefreshToken: "invalid-refresh-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), "invalid-refresh-token").
					Return(nil, &security.InvalidJWTError{Message: "invalid token"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Unauthenticated, Message: "invalid token"},
			errorExpected: true,
		},
		{
			name: "access token mismatch",
			in:   &sso.RefreshTokensIn{RefreshToken: "valid-refresh-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), "valid-refresh-token").
					Return(nil, &customerrors.AccessTokenDoesNotBelongToRefreshTokenError{Message: "token mismatch"}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErr:   &customgrpc.BaseError{Status: codes.Unauthenticated, Message: "token mismatch"},
			errorExpected: true,
		},
		{
			name: "user not found",
			in:   &sso.RefreshTokensIn{RefreshToken: "valid-refresh-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), "valid-refresh-token").
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
			in:   &sso.RefreshTokensIn{RefreshToken: "valid-refresh-token"},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogging.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), "valid-refresh-token").
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

			resp, err := api.RefreshTokens(context.Background(), tc.in)
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
