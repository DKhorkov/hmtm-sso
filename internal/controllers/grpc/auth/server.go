package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/internal/entities"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

// Register handler registers new User with provided data.
func (api *ServerAPI) Register(ctx context.Context, request *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	userData := entities.RegisterUserDTO{
		Credentials: entities.LoginUserDTO{
			Email:    request.GetEmail(),
			Password: request.GetPassword(),
		},
	}

	userID, err := api.useCases.RegisterUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to register User", err)

		switch {
		case errors.As(err, &customerrors.UserAlreadyExistsError{}):
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.RegisterResponse{UserID: userID}, nil
}

// Login handler authenticates user if provided credentials are valid and logs User in system.
func (api *ServerAPI) Login(ctx context.Context, request *sso.LoginRequest) (*sso.LoginResponse, error) {
	userData := entities.LoginUserDTO{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
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
		case errors.As(err, &customerrors.InvalidPasswordError{}):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.LoginResponse{
		AccessToken:  tokensDTO.AccessToken,
		RefreshToken: tokensDTO.RefreshToken,
	}, nil
}

// RefreshTokens handler updates User auth tokens.
func (api *ServerAPI) RefreshTokens(
	ctx context.Context,
	request *sso.RefreshTokensRequest,
) (*sso.LoginResponse, error) {
	tokensDTO, err := api.useCases.RefreshTokens(ctx, request.GetRefreshToken())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to refresh tokens with RefreshToken=%s", request.GetRefreshToken()),
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

	return &sso.LoginResponse{
		AccessToken:  tokensDTO.AccessToken,
		RefreshToken: tokensDTO.RefreshToken,
	}, nil
}
