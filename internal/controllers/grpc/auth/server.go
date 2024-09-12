package auth

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	// Helps to test single endpoints, if others is not implemented yet
	sso.UnimplementedAuthServiceServer
}

// Register handler (serverAPI) for AuthServer  to gRPC server:.
func Register(gRPCServer *grpc.Server) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{})
}

// Register user handler for AuthServer.
func (api *ServerAPI) Register(ctx context.Context, request *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	return &sso.RegisterResponse{UserID: 1}, nil
}

// Login user handler for AuthServer.
func (api *ServerAPI) Login(ctx context.Context, request *sso.LoginRequest) (*sso.LoginResponse, error) {
	return &sso.LoginResponse{Token: "someToken"}, nil
}
