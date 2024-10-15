package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"

	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
	"google.golang.org/grpc/codes"

	"github.com/DKhorkov/hmtm-sso/pkg/entities"

	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"google.golang.org/grpc"

	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedAuthServiceServer
	UseCases interfaces.UseCases
	Logger   *slog.Logger
}

// Register handler registers new User with provided data.
func (api *ServerAPI) Register(ctx context.Context, request *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	api.Logger.InfoContext(
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

	userID, err := api.UseCases.RegisterUser(userData)
	if err != nil {
		api.Logger.ErrorContext(
			ctx,
			"Error occurred while trying to register",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var userAlreadyExists *customerrors.UserAlreadyExistsError
		if errors.As(err, &userAlreadyExists) {
			return nil, &customerrors.GRPCError{Status: codes.AlreadyExists, Message: err.Error()}
		}

		return nil, &customerrors.GRPCError{Status: codes.Internal, Message: err.Error()}
	}

	return &sso.RegisterResponse{UserID: int64(userID)}, nil
}

// Login handler authenticates user if provided credentials are valid and logs User in system.
func (api *ServerAPI) Login(ctx context.Context, request *sso.LoginRequest) (*sso.LoginResponse, error) {
	api.Logger.InfoContext(
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

	tokensDTO, err := api.UseCases.LoginUser(userData)
	if err != nil {
		api.Logger.ErrorContext(
			ctx,
			"Error occurred while trying to login",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var userNotFoundError *customerrors.UserNotFoundError
		if errors.As(err, &userNotFoundError) {
			return nil, &customerrors.GRPCError{Status: codes.NotFound, Message: err.Error()}
		}

		var invalidPasswordError *customerrors.InvalidPasswordError
		if errors.As(err, &invalidPasswordError) {
			return nil, &customerrors.GRPCError{Status: codes.Unauthenticated, Message: err.Error()}
		}

		return nil, &customerrors.GRPCError{Status: codes.Internal, Message: err.Error()}
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
	api.Logger.InfoContext(
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

	tokensDTO, err := api.UseCases.RefreshTokens(refreshTokensDTO)
	if err != nil {
		api.Logger.ErrorContext(
			ctx,
			"Error occurred while trying to login",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var invalidJWTError *customerrors.InvalidJWTError
		var accessTokenDoesNotBelongToRefreshTokenError *customerrors.AccessTokenDoesNotBelongToRefreshTokenError
		if errors.As(err, &invalidJWTError) || errors.As(err, &accessTokenDoesNotBelongToRefreshTokenError) {
			return nil, &customerrors.GRPCError{Status: codes.Unauthenticated, Message: err.Error()}
		}

		var userNotFoundError *customerrors.UserNotFoundError
		if errors.As(err, &userNotFoundError) {
			return nil, &customerrors.GRPCError{Status: codes.NotFound, Message: err.Error()}
		}

		return nil, &customerrors.GRPCError{Status: codes.Internal, Message: err.Error()}
	}

	return &sso.LoginResponse{
		AccessToken:  tokensDTO.AccessToken,
		RefreshToken: tokensDTO.RefreshToken,
	}, nil
}

// RegisterServer handler (serverAPI) for AuthServer  to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{UseCases: useCases, Logger: logger})
}
