package auth

import (
	"context"
	"errors"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"github.com/DKhorkov/libs/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	customgrpc "github.com/DKhorkov/libs/grpc"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

var (
	userNotFoundError                           = &customerrors.UserNotFoundError{}
	userAlreadyExistsError                      = &customerrors.UserAlreadyExistsError{}
	emailAlreadyConfirmedError                  = &customerrors.EmailAlreadyConfirmedError{}
	emailIsNotConfirmedError                    = &customerrors.EmailIsNotConfirmedError{}
	invalidJWTError                             = &security.InvalidJWTError{}
	wrongPasswordError                          = &customerrors.WrongPasswordError{}
	accessTokenDoesNotBelongToRefreshTokenError = &customerrors.AccessTokenDoesNotBelongToRefreshTokenError{}
	validationError                             = &validation.Error{}
)

// RegisterServer handler (serverAPI) for AuthServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger logging.Logger) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedAuthServiceServer
	useCases interfaces.UseCases
	logger   logging.Logger
}

func (api *ServerAPI) SendForgetPasswordMessage(ctx context.Context, in *sso.SendForgetPasswordMessageIn) (*emptypb.Empty, error) {
	if err := api.useCases.SendForgetPasswordMessage(ctx, in.GetEmail()); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to send forget-password message to User with email="+in.GetEmail(),
			err,
		)

		switch {
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &emailIsNotConfirmedError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) SendVerifyEmailMessage(
	ctx context.Context,
	in *sso.SendVerifyEmailMessageIn,
) (*emptypb.Empty, error) {
	if err := api.useCases.SendVerifyEmailMessage(ctx, in.GetEmail()); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to send verify-email message to User with email="+in.GetEmail(),
			err,
		)

		switch {
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &emailAlreadyConfirmedError):
			return nil, &customgrpc.BaseError{
				Status:  codes.FailedPrecondition,
				Message: err.Error(),
			}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) ChangePassword(
	ctx context.Context,
	in *sso.ChangePasswordIn,
) (*emptypb.Empty, error) {
	if err := api.useCases.ChangePassword(
		ctx,
		in.GetAccessToken(),
		in.GetOldPassword(),
		in.GetNewPassword(),
	); err != nil {
		switch {
		case errors.As(err, &validationError),
			errors.As(err, &wrongPasswordError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) ForgetPassword(
	ctx context.Context,
	in *sso.ForgetPasswordIn,
) (*emptypb.Empty, error) {
	if err := api.useCases.ForgetPassword(ctx, in.GetForgetPasswordToken(), in.GetNewPassword()); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to change forgotten password for User with forgetPasswordToken="+in.GetForgetPasswordToken(),
			err,
		)

		switch {
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &validationError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) VerifyEmail(
	ctx context.Context,
	in *sso.VerifyEmailIn,
) (*emptypb.Empty, error) {
	if err := api.useCases.VerifyUserEmail(ctx, in.GetVerifyEmailToken()); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to verify email for User with VerifyEmailToken="+in.GetVerifyEmailToken(),
			err,
		)

		switch {
		case errors.As(err, &emailAlreadyConfirmedError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) Logout(ctx context.Context, in *sso.LogoutIn) (*emptypb.Empty, error) {
	if err := api.useCases.LogoutUser(ctx, in.GetAccessToken()); err != nil {
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	return &emptypb.Empty{}, nil
}

// Register handler registers new User with provided data.
func (api *ServerAPI) Register(ctx context.Context, in *sso.RegisterIn) (*sso.RegisterOut, error) {
	userData := entities.RegisterUserDTO{
		DisplayName: in.GetDisplayName(),
		Email:       in.GetEmail(),
		Password:    in.GetPassword(),
	}

	userID, err := api.useCases.RegisterUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to register User",
			err,
		)

		switch {
		case errors.As(err, &validationError):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		case errors.As(err, &userAlreadyExistsError):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.RegisterOut{UserID: userID}, nil
}

// Login handler authenticates User if provided credentials are valid and logs User in system.
func (api *ServerAPI) Login(ctx context.Context, in *sso.LoginIn) (*sso.LoginOut, error) {
	userData := entities.LoginUserDTO{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	tokensDTO, err := api.useCases.LoginUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to login User with email="+userData.Email,
			err,
		)

		switch {
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &wrongPasswordError):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.LoginOut{
		AccessToken:  tokensDTO.AccessToken,
		RefreshToken: tokensDTO.RefreshToken,
	}, nil
}

// RefreshTokens handler updates User auth tokens.
func (api *ServerAPI) RefreshTokens(
	ctx context.Context,
	in *sso.RefreshTokensIn,
) (*sso.LoginOut, error) {
	tokensDTO, err := api.useCases.RefreshTokens(ctx, in.GetRefreshToken())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to refresh tokens with RefreshToken="+in.GetRefreshToken(),
			err,
		)

		switch {
		case errors.As(err, &invalidJWTError),
			errors.As(err, &accessTokenDoesNotBelongToRefreshTokenError):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.LoginOut{
		AccessToken:  tokensDTO.AccessToken,
		RefreshToken: tokensDTO.RefreshToken,
	}, nil
}
