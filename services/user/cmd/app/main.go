package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/guluzadehh/go_bookstore/services/user/internal/app"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	"github.com/joho/godotenv"
)

const (
	env_local = "local"
	env_dev   = "dev"
	env_prod  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %s\n", err.Error())
	}

	config := config.MustLoad()

	log := setupLogger(config.Env)

	log.Info("building user app", slog.String("env", config.Env))
	app := app.New(log, config)

	log.Info("starting user app")
	go app.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("stopping user app")
	app.Stop()

	log.Info("user app has been gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case env_local, env_dev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env_prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
