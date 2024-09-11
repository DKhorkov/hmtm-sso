package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/DKhorkov/hmtm-sso/internal/grpc/auth"
	usersgrpc "github.com/DKhorkov/hmtm-sso/internal/grpc/users"
	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"google.golang.org/grpc"
)

type GrpcApp struct {
	grpcServer *grpc.Server
	host       string
	port       int
	logger     *slog.Logger
}

// Run gRPC server.
func (grpcApp *GrpcApp) Run() {
	grpcApp.logger.Info(
		fmt.Sprintf("Starting gRPC Server at %s:%d", grpcApp.host, grpcApp.port),
		"Traceback",
		logging.GetLogTraceback(),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcApp.host, grpcApp.port))
	if err != nil {
		grpcApp.logger.Error(
			"Failed to start gRPC Server",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
		panic(err)
	}

	if err = grpcApp.grpcServer.Serve(listener); err != nil {
		grpcApp.logger.Error(
			"Error occurred while listening to gRPC server",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
		panic(err)
	}
}

// Stop gRPC server gracefully (graceful shutdown).
func (grpcApp *GrpcApp) Stop() {
	// Stops accepting new requests and processes already received requests:
	grpcApp.grpcServer.GracefulStop()
}

// New creates an instance of GrpcApp like a constructor.
func New(host string, port int) *GrpcApp {
	grpcServer := grpc.NewServer()

	// Connects our gRPC services to grpcServer:
	authgrpc.Register(grpcServer)
	usersgrpc.Register(grpcServer)

	return &GrpcApp{
		grpcServer: grpcServer,
		port:       port,
		host:       host,
		logger:     logging.GetInstance(logging.LogLevels.DEBUG),
	}
}
