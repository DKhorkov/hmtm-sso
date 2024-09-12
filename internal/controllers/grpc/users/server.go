package users

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedUsersServiceServer
}

// Register handler (serverAPI) for AuthServer  to gRPC server:.
func Register(gRPCServer *grpc.Server) {
	sso.RegisterUsersServiceServer(gRPCServer, &ServerAPI{})
}

// GetUser handler return user with provided data for UsersServer.
func (api *ServerAPI) GetUser(ctx context.Context, request *sso.GetUserRequest) (*sso.GetUserResponse, error) {
	response := &sso.GetUserResponse{
		UserID:    request.GetUserID(),
		Email:     "SOmeEmail",
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	return response, nil
}

// GetUsers user handler return all users for UsersServer.
func (api *ServerAPI) GetUsers(ctx context.Context, request *emptypb.Empty) (*sso.GetUsersResponse, error) {
	response := &sso.GetUsersResponse{
		Users: []*sso.GetUserResponse{
			{
				UserID:    1,
				Email:     "SOmeEmail",
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
		},
	}

	return response, nil
}
