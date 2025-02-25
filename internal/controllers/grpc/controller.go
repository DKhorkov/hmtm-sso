package grpccontroller

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	customgrpc "github.com/DKhorkov/libs/grpc/interceptors"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-sso/internal/controllers/grpc/auth"
	"github.com/DKhorkov/hmtm-sso/internal/controllers/grpc/users"
	"github.com/DKhorkov/hmtm-sso/internal/interfaces"
)

// New creates an instance of gRPC Controller.
func New(
	host string,
	port int,
	useCases interfaces.UseCases,
	logger logging.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *Controller {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			customgrpc.UnaryServerTracingInterceptor(traceProvider, spanConfig),
			customgrpc.UnaryServerLoggingInterceptor(logger),
		),
	)

	// Connects our gRPC services to grpcServer:
	auth.RegisterServer(grpcServer, useCases, logger)
	users.RegisterServer(grpcServer, useCases, logger)

	return &Controller{
		grpcServer: grpcServer,
		port:       port,
		host:       host,
		logger:     logger,
	}
}

type Controller struct {
	grpcServer *grpc.Server
	host       string
	port       int
	logger     logging.Logger
}

// Run gRPC server.
func (controller *Controller) Run() {
	logging.LogInfo(
		controller.logger,
		fmt.Sprintf("Starting gRPC Server at http://%s:%d", controller.host, controller.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", controller.host, controller.port))
	if err != nil {
		logging.LogError(controller.logger, "Failed to start gRPC Server", err)
		panic(err)
	}

	if err = controller.grpcServer.Serve(listener); err != nil {
		logging.LogError(controller.logger, "Error occurred while listening to gRPC server", err)
		panic(err)
	}

	logging.LogInfo(controller.logger, "Stopped serving new connections.")
}

// Stop gRPC server gracefully (graceful shutdown).
func (controller *Controller) Stop() {
	// Stops accepting new requests and processes already received requests:
	controller.grpcServer.GracefulStop()
	logging.LogInfo(controller.logger, "Graceful shutdown completed.")
}
