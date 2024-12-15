package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/requestid"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	customerrors "github.com/DKhorkov/hmtm-sso/internal/errors"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	user, err := api.useCases.GetUserByID(ctx, request.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get User with ID=%d", request.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.UserNotFoundError{}):
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
func (api *ServerAPI) GetUsers(ctx context.Context, request *sso.GetUsersRequest) (*sso.GetUsersResponse, error) {
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	users, err := api.useCases.GetAllUsers(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Users", err)
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
	ctx = contextlib.SetValue(ctx, requestid.Key, request.GetRequestID())
	logging.LogRequest(ctx, api.logger, request)

	user, err := api.useCases.GetMe(ctx, request.GetAccessToken())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get User with AccessToken=%s", request.GetAccessToken()),
			err,
		)

		switch {
		case errors.As(err, &security.InvalidJWTError{}):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &customerrors.UserNotFoundError{}):
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
