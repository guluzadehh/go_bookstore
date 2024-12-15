package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	usergrpc "github.com/guluzadehh/go_bookstore/services/user/internal/grpc/user"
	"google.golang.org/grpc"
)

type GrpcApp struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	config *config.Config,
	userService usergrpc.UserService,
) *GrpcApp {
	server := grpc.NewServer()

	userHandler := usergrpc.New(log, config, userService)

	usergrpc.Register(server, userHandler)

	return &GrpcApp{
		log:        log,
		grpcServer: server,
		port:       config.GRPCServer.Port,
	}
}

func (a *GrpcApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *GrpcApp) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting user grpc server", slog.String("addr", l.Addr().String()))

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *GrpcApp) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(slog.String("op", op))

	a.grpcServer.GracefulStop()
	log.Info("user grpc server has been gracefully stopped")
}
