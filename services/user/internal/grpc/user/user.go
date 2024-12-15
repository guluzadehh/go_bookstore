package usergrpc

import (
	"log/slog"

	user_grpc "github.com/guluzadehh/go_bookstore/protos/gen/go/user"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	grpchandler "github.com/guluzadehh/go_bookstore/services/user/internal/grpc"
	"google.golang.org/grpc"
)

type UserService interface{}

type UserHandler struct {
	*grpchandler.Handler
	cfg  *config.Config
	srvc UserService
	user_grpc.UnimplementedUserServiceServer
}

func New(log *slog.Logger, cfg *config.Config, userService UserService) *UserHandler {
	return &UserHandler{
		Handler: grpchandler.NewHandler(log),
		cfg:     cfg,
		srvc:    userService,
	}
}

func Register(grpcServer *grpc.Server, userHandler *UserHandler) {
	user_grpc.RegisterUserServiceServer(grpcServer, userHandler)
}
