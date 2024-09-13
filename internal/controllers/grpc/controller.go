package grpccontroller

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/DKhorkov/hmtm-sso/internal/controllers/grpc/auth"
	"github.com/DKhorkov/hmtm-sso/internal/controllers/grpc/users"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"google.golang.org/grpc"
)

type Controller struct {
	grpcServer *grpc.Server
	host       string
	port       int
	logger     *slog.Logger
}

// Run gRPC server.
func (controller *Controller) Run() {
	controller.logger.Info(
		fmt.Sprintf("Starting gRPC Server at %s:%d", controller.host, controller.port),
		"Traceback",
		logging.GetLogTraceback(),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", controller.host, controller.port))
	if err != nil {
		controller.logger.Error(
			"Failed to start gRPC Server",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
		panic(err)
	}

	if err = controller.grpcServer.Serve(listener); err != nil {
		controller.logger.Error(
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
func (controller *Controller) Stop() {
	// Stops accepting new requests and processes already received requests:
	controller.grpcServer.GracefulStop()
}

// New creates an instance of GrpcApp like a constructor.
func New(host string, port int) *Controller {
	grpcServer := grpc.NewServer()

	// Connects our gRPC services to grpcServer:
	auth.Register(grpcServer)
	users.Register(grpcServer)

	return &Controller{
		grpcServer: grpcServer,
		port:       port,
		host:       host,
		logger:     logging.GetInstance(logging.LogLevels.DEBUG),
	}
}
