package auth

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
	"google.golang.org/grpc"

	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedAuthServiceServer
	UseCases interfaces.UseCases
}

// Register user handler for AuthServer.
func (api *ServerAPI) Register(ctx context.Context, request *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	userData := entities.RegisterUserDTO{
		Credentials: entities.LoginUserDTO{
			Email:    request.GetCredentials().GetEmail(),
			Password: request.GetCredentials().GetPassword(),
		},
	}

	userID, err := api.UseCases.RegisterUser(userData)
	if err != nil {
		return nil, err
	}

	return &sso.RegisterResponse{UserID: int64(userID)}, nil
}

// Login user handler for AuthServer.
func (api *ServerAPI) Login(ctx context.Context, request *sso.LoginRequest) (*sso.LoginResponse, error) {
	userData := entities.LoginUserDTO{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	token, err := api.UseCases.LoginUser(userData)
	if err != nil {
		return &sso.LoginResponse{Token: ""}, err
	}

	return &sso.LoginResponse{Token: token}, nil
}

// Register handler (serverAPI) for AuthServer  to gRPC server:.
func Register(gRPCServer *grpc.Server, useCases interfaces.UseCases) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{UseCases: useCases})
}
