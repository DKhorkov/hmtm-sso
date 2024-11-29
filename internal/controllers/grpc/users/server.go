package users

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedUsersServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

// GetUser handler returns User according provided data.
func (api *ServerAPI) GetUser(ctx context.Context, request *sso.GetUserRequest) (*sso.GetUserResponse, error) {
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

	user, err := api.useCases.GetUserByID(request.GetID())
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get user",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var userNotFoundError *customerrors.UserNotFoundError
		switch {
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// GetUsers handler returns all Users.
func (api *ServerAPI) GetUsers(ctx context.Context, request *emptypb.Empty) (*sso.GetUsersResponse, error) {
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

	users, err := api.useCases.GetAllUsers()
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get all users",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	usersForResponse := make([]*sso.GetUserResponse, len(users))
	for i, user := range users {
		usersForResponse[i] = &sso.GetUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		}
	}

	return &sso.GetUsersResponse{Users: usersForResponse}, nil
}

// GetMe handler returns User according to provided Access Token.
func (api *ServerAPI) GetMe(ctx context.Context, request *sso.GetMeRequest) (*sso.GetUserResponse, error) {
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

	user, err := api.useCases.GetMe(request.GetAccessToken())
	if err != nil {
		api.logger.ErrorContext(
			ctx,
			"Error occurred while trying to get user",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var invalidJWTError *security.InvalidJWTError
		var userNotFoundError *customerrors.UserNotFoundError
		switch {
		case errors.As(err, &invalidJWTError):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &sso.GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// RegisterServer handler (serverAPI) for UsersServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	sso.RegisterUsersServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}
