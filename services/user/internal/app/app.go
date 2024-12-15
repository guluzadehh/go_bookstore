package app

import (
	"log/slog"

	grpcapp "github.com/guluzadehh/go_bookstore/services/user/internal/app/grpc"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	"github.com/guluzadehh/go_bookstore/services/user/internal/services/user"
)

type App struct {
	log     *slog.Logger
	grpcApp *grpcapp.GrpcApp
}

func New(log *slog.Logger, config *config.Config) *App {
	userService := user.New(log, config)

	grpcApp := grpcapp.New(log, config, userService)

	return &App{
		log:     log,
		grpcApp: grpcApp,
	}
}

func (a *App) Start() {
	a.log.Info("starting user grpc app")
	a.grpcApp.MustRun()
}

func (a *App) Stop() {
	a.log.Info("stopping user grpc app")
	a.grpcApp.Stop()
}
