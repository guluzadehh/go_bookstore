package auth

import (
	"log/slog"

	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
)


type AuthService struct {
	log          *slog.Logger
	config       *config.Config
}

func New(log *slog.Logger, config *config.Config) *AuthService {
	return &AuthService{
		log:	log,
		config:	config,
	}
}
}
