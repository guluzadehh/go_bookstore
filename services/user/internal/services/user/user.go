package user

import (
	"log/slog"

	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
)

type UserService struct {
	log    *slog.Logger
	config *config.Config
}

func New(log *slog.Logger, config *config.Config) *UserService {
	return &UserService{
		log:    log,
		config: config,
	}
}
