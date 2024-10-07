package users

import (
	"context"
	"errors"

	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	customerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedUsersServiceServer
	UseCases interfaces.UseCases
}

// GetUser handler return user with provided data for UsersServer.
func (api *ServerAPI) GetUser(ctx context.Context, request *sso.GetUserRequest) (*sso.GetUserResponse, error) {
	user, err := api.UseCases.GetUserByID(int(request.GetUserID()))
	if err != nil {
		var userNotFoundError *customerrors.UserNotFoundError
		if errors.As(err, &userNotFoundError) {
			return nil, &customerrors.GRPCError{Status: codes.NotFound, Message: err.Error()}
		}

		return nil, &customerrors.GRPCError{Status: codes.Internal, Message: err.Error()}
	}

	response := &sso.GetUserResponse{
		UserID:    int64(user.ID),
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	return response, nil
}

// GetUsers user handler return all users for UsersServer.
func (api *ServerAPI) GetUsers(ctx context.Context, request *emptypb.Empty) (*sso.GetUsersResponse, error) {
	users, err := api.UseCases.GetAllUsers()
	if err != nil {
		return nil, &customerrors.GRPCError{Status: codes.Internal, Message: err.Error()}
	}

	usersForResponse := make([]*sso.GetUserResponse, len(users))
	for i, user := range users {
		usersForResponse[i] = &sso.GetUserResponse{
			UserID:    int64(user.ID),
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		}
	}

	response := &sso.GetUsersResponse{
		Users: usersForResponse,
	}

	return response, nil
}

// Register handler (serverAPI) for AuthServer  to gRPC server:.
func Register(gRPCServer *grpc.Server, useCases interfaces.UseCases) {
	sso.RegisterUsersServiceServer(gRPCServer, &ServerAPI{UseCases: useCases})
}
