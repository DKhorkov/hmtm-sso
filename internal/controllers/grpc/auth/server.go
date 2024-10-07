package auth

import (
	"context"
	"errors"

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
		var userAlreadyExists *customerrors.UserAlreadyExistsError
		if errors.As(err, &userAlreadyExists) {
			return nil, &customerrors.GRPCError{Status: codes.AlreadyExists, Message: err.Error()}
		}

		return nil, &customerrors.GRPCError{Status: codes.Internal, Message: err.Error()}
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

	return &sso.LoginResponse{Token: token}, nil
}

// Register handler (serverAPI) for AuthServer  to gRPC server:.
func Register(gRPCServer *grpc.Server, useCases interfaces.UseCases) {
	sso.RegisterAuthServiceServer(gRPCServer, &ServerAPI{UseCases: useCases})
}
