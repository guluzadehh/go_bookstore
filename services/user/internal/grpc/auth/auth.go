package authgrpc

import (
	"context"
	"errors"
	"log/slog"

	user_grpc "github.com/guluzadehh/go_bookstore/protos/gen/go/user"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	grpchandler "github.com/guluzadehh/go_bookstore/services/user/internal/grpc"
	service "github.com/guluzadehh/go_bookstore/services/user/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Authenticate(ctx context.Context, email string, password string) (access string, refresh string, err error)
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

func (h *AuthHandler) Authenticate(ctx context.Context, req *user_grpc.AuthRequest) (*user_grpc.AuthResponse, error) {
	access, refresh, err := h.srvc.Authenticate(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "failed to authenticate")
	}

	return &user_grpc.AuthResponse{
		Access:  access,
		Refresh: refresh,
	}, nil
}
