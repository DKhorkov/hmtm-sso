package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/pkg/entities"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedAuthServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// Register handler registers new User with provided data.
func (api *ServerAPI) Register(ctx context.Context, request *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	api.logger.InfoContext(
		ctx,
		"Received new request",
		"Request",
		request,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userData := entities.RegisterUserDTO{
		Credentials: entities.LoginUserDTO{
			Email:    request.GetCredentials().GetEmail(),
			Password: request.GetCredentials().GetPassword(),
		},
	}

	userID, err := api.useCases.RegisterUser(userData)
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to register",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var userAlreadyExists *customerrors.UserAlreadyExistsError
		if errors.As(err, &userAlreadyExists) {
			return nil, &customgrpc.BaseError{Status: codes.AlreadyExists, Message: err.Error()}
		}

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	return &sso.RegisterResponse{UserID: userID}, nil
}

// Login handler authenticates user if provided credentials are valid and logs User in system.
func (api *ServerAPI) Login(ctx context.Context, request *sso.LoginRequest) (*sso.LoginResponse, error) {
	api.logger.InfoContext(
		ctx,
		"Received new request",
		"Request",
		request,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userData := entities.LoginUserDTO{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	tokensDTO, err := api.useCases.LoginUser(userData)
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to login",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var userNotFoundError *customerrors.UserNotFoundError
		if errors.As(err, &userNotFoundError) {
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		}

		var invalidPasswordError *customerrors.InvalidPasswordError
		if errors.As(err, &invalidPasswordError) {
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		}

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
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
	api.logger.InfoContext(
		ctx,
		"Received new request",
		"Request",
		request,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	refreshTokensDTO := entities.TokensDTO{
		AccessToken:  request.GetAccessToken(),
		RefreshToken: request.GetRefreshToken(),
	}

	tokensDTO, err := api.useCases.RefreshTokens(refreshTokensDTO)
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to login",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var invalidJWTError *security.InvalidJWTError
		var accessTokenDoesNotBelongToRefreshTokenError *customerrors.AccessTokenDoesNotBelongToRefreshTokenError
		if errors.As(err, &invalidJWTError) || errors.As(err, &accessTokenDoesNotBelongToRefreshTokenError) {
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		}

		var userNotFoundError *customerrors.UserNotFoundError
		if errors.As(err, &userNotFoundError) {
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		}

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	return &sso.LoginResponse{
		AccessToken:  tokensDTO.AccessToken,
		RefreshToken: tokensDTO.RefreshToken,
	}, nil
}

// RegisterServer handler (serverAPI) for AuthServer  to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}
