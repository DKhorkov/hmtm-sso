package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

// RegisterServer handler (serverAPI) for AuthServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedAuthServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

func (api *ServerAPI) SendVerifyEmailMessage(
	ctx context.Context,
	in *sso.SendVerifyEmailMessageIn,
) (*emptypb.Empty, error) {
	if err := api.useCases.SendVerifyEmailMessage(ctx, in.GetEmail()); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to login User with email=%s", in.GetEmail()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.UserNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &customerrors.EmailAlreadyConfirmedError{}):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) ChangePassword(ctx context.Context, in *sso.ChangePasswordIn) (*emptypb.Empty, error) {
	if err := api.useCases.ChangePassword(
		ctx,
		in.GetAccessToken(),
		in.GetOldPassword(),
		in.GetNewPassword(),
	); err != nil {
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) ForgetPassword(ctx context.Context, in *sso.ForgetPasswordIn) (*emptypb.Empty, error) {
	if err := api.useCases.ForgetPassword(ctx, in.GetAccessToken()); err != nil {
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) VerifyEmail(ctx context.Context, in *sso.VerifyEmailIn) (*emptypb.Empty, error) {
	if err := api.useCases.VerifyUserEmail(ctx, in.GetVerifyEmailToken()); err != nil {
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
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
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to register User", err)

		switch {
		case errors.As(err, &customerrors.InvalidEmailError{}),
			errors.As(err, &customerrors.InvalidPasswordError{}):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		case errors.As(err, &customerrors.UserAlreadyExistsError{}):
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
			fmt.Sprintf("Error occurred while trying to login User with email=%s", userData.Email),
			err,
		)

		switch {
		case errors.As(err, &customerrors.UserNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		case errors.As(err, &customerrors.WrongPasswordError{}):
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
func (api *ServerAPI) RefreshTokens(ctx context.Context, in *sso.RefreshTokensIn) (*sso.LoginOut, error) {
	tokensDTO, err := api.useCases.RefreshTokens(ctx, in.GetRefreshToken())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to refresh tokens with RefreshToken=%s", in.GetRefreshToken()),
			err,
		)

		switch {
		case errors.As(err, &security.InvalidJWTError{}),
			errors.As(err, &customerrors.AccessTokenDoesNotBelongToRefreshTokenError{}):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &customerrors.UserNotFoundError{}):
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
