package users

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

// RegisterServer handler (serverAPI) for UsersServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger *slog.Logger) {
	sso.RegisterUsersServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedUsersServiceServer
	useCases interfaces.UseCases
	logger   *slog.Logger
}

func (api *ServerAPI) UpdateUserProfile(ctx context.Context, in *sso.UpdateUserProfileIn) (*emptypb.Empty, error) {
	userProfileData := entities.RawUpdateUserProfileDTO{
		AccessToken: in.GetAccessToken(),
		DisplayName: in.GetDisplayName(),
		Phone:       in.GetPhone(),
		Telegram:    in.GetTelegram(),
		Avatar:      in.GetAvatar(),
	}

	if err := api.useCases.UpdateUserProfile(ctx, userProfileData); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to update User profile with AccessToken=%s", in.GetAccessToken()),
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

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) GetUserByEmail(ctx context.Context, in *sso.GetUserByEmailIn) (*sso.GetUserOut, error) {
	user, err := api.useCases.GetUserByEmail(ctx, in.GetEmail())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get User with Email=%s", in.GetEmail()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.UserNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return prepareUserOut(*user), nil
}

// GetUser handler returns User according provided data.
func (api *ServerAPI) GetUser(ctx context.Context, in *sso.GetUserIn) (*sso.GetUserOut, error) {
	user, err := api.useCases.GetUserByID(ctx, in.GetID())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get User with ID=%d", in.GetID()),
			err,
		)

		switch {
		case errors.As(err, &customerrors.UserNotFoundError{}):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return prepareUserOut(*user), nil
}

// GetUsers handler returns all Users.
func (api *ServerAPI) GetUsers(ctx context.Context, _ *emptypb.Empty) (*sso.GetUsersOut, error) {
	users, err := api.useCases.GetAllUsers(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, api.logger, "Error occurred while trying to get all Users", err)
		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedUsers := make([]*sso.GetUserOut, len(users))
	for i, user := range users {
		processedUsers[i] = prepareUserOut(user)
	}

	return &sso.GetUsersOut{Users: processedUsers}, nil
}

// GetMe handler returns User according to provided Access Token.
func (api *ServerAPI) GetMe(ctx context.Context, in *sso.GetMeIn) (*sso.GetUserOut, error) {
	user, err := api.useCases.GetMe(ctx, in.GetAccessToken())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			fmt.Sprintf("Error occurred while trying to get User with AccessToken=%s", in.GetAccessToken()),
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

	return prepareUserOut(*user), nil
}
