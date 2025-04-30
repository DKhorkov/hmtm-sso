package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/security"
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
	userNotFoundError  = &customerrors.UserNotFoundError{}
	invalidJWTError    = &security.InvalidJWTError{}
	invalidDisplayName = &customerrors.InvalidDisplayNameError{}
	invalidPhone       = &customerrors.InvalidPhoneError{}
	invalidTelegramErr = &customerrors.InvalidTelegramError{}
)

// RegisterServer handler (serverAPI) for UsersServer to gRPC server:.
func RegisterServer(gRPCServer *grpc.Server, useCases interfaces.UseCases, logger logging.Logger) {
	sso.RegisterUsersServiceServer(gRPCServer, &ServerAPI{useCases: useCases, logger: logger})
}

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedUsersServiceServer
	useCases interfaces.UseCases
	logger   logging.Logger
}

func (api *ServerAPI) UpdateUserProfile(
	ctx context.Context,
	in *sso.UpdateUserProfileIn,
) (*emptypb.Empty, error) {
	userProfileData := entities.RawUpdateUserProfileDTO{
		AccessToken: in.GetAccessToken(),
		DisplayName: in.DisplayName,
		Phone:       in.Phone,
		Telegram:    in.Telegram,
		Avatar:      in.Avatar,
	}

	if err := api.useCases.UpdateUserProfile(ctx, userProfileData); err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to update User profile with AccessToken="+in.GetAccessToken(),
			err,
		)

		switch {
		case errors.As(err, &invalidPhone),
			errors.As(err, &invalidTelegramErr),
			errors.As(err, &invalidDisplayName):
			return nil, &customgrpc.BaseError{Status: codes.FailedPrecondition, Message: err.Error()}
		case errors.As(err, &invalidJWTError):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return &emptypb.Empty{}, nil
}

func (api *ServerAPI) GetUserByEmail(
	ctx context.Context,
	in *sso.GetUserByEmailIn,
) (*sso.GetUserOut, error) {
	user, err := api.useCases.GetUserByEmail(ctx, in.GetEmail())
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to get User with Email="+in.GetEmail(),
			err,
		)

		switch {
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return mapUserToOut(*user), nil
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
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return mapUserToOut(*user), nil
}

// GetUsers handler returns all Users.
func (api *ServerAPI) GetUsers(ctx context.Context, in *sso.GetUsersIn) (*sso.GetUsersOut, error) {
	var pagination *entities.Pagination
	if in.Pagination != nil {
		pagination = &entities.Pagination{
			Limit:  in.Pagination.Limit,
			Offset: in.Pagination.Offset,
		}
	}

	users, err := api.useCases.GetUsers(ctx, pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			api.logger,
			"Error occurred while trying to get all Users",
			err,
		)

		return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
	}

	processedUsers := make([]*sso.GetUserOut, len(users))
	for i, user := range users {
		processedUsers[i] = mapUserToOut(user)
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
			"Error occurred while trying to get User with AccessToken="+in.GetAccessToken(),
			err,
		)

		switch {
		case errors.As(err, &invalidJWTError):
			return nil, &customgrpc.BaseError{Status: codes.Unauthenticated, Message: err.Error()}
		case errors.As(err, &userNotFoundError):
			return nil, &customgrpc.BaseError{Status: codes.NotFound, Message: err.Error()}
		default:
			return nil, &customgrpc.BaseError{Status: codes.Internal, Message: err.Error()}
		}
	}

	return mapUserToOut(*user), nil
}
