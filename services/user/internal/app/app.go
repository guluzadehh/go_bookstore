package app

import (
	"log/slog"

	grpcapp "github.com/guluzadehh/go_bookstore/services/user/internal/app/grpc"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	"github.com/guluzadehh/go_bookstore/services/user/internal/storage/postgresql"
)

type App struct {
	log     *slog.Logger
	grpcApp *grpcapp.GrpcApp
}

func New(log *slog.Logger, config *config.Config) *App {
	grpcApp := grpcapp.New(log, config)

	_, err := postgresql.New(config.Postgresql.DSN(nil))
	if err != nil {
		panic(err)
	}

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
