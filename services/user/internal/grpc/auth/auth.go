package authgrpc

import (
	"log/slog"

	user_grpc "github.com/guluzadehh/go_bookstore/protos/gen/go/user"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	grpchandler "github.com/guluzadehh/go_bookstore/services/user/internal/grpc"
)

type AuthService interface {
}

type AuthHandler struct {
	*grpchandler.Handler
	cfg  *config.Config
	srvc AuthService
	user_grpc.UnimplementedAuthServiceServer
}

func New(log *slog.Logger, cfg *config.Config, srvc AuthService) *AuthHandler {
	return &AuthHandler{
		Handler: grpchandler.NewHandler(log),
		cfg:     cfg,
		srvc:    srvc,
	}
}

func Register(grpcServer *grpc.Server, authHandler *AuthHandler) {
	user_grpc.RegisterAuthServiceServer(grpcServer, authHandler)
}

}
