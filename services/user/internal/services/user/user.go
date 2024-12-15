package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	"github.com/guluzadehh/go_bookstore/services/user/internal/domain/models"
	"github.com/guluzadehh/go_bookstore/services/user/internal/lib/sl"
	service "github.com/guluzadehh/go_bookstore/services/user/internal/services"
	"github.com/guluzadehh/go_bookstore/services/user/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	CreateUser(ctx context.Context, email string, password string) (*models.User, error)
}

type UserService struct {
	log       *slog.Logger
	config    *config.Config
	userSaver UserSaver
}

func New(log *slog.Logger, config *config.Config, userSaver UserSaver) *UserService {
	return &UserService{
		log:       log,
		config:    config,
		userSaver: userSaver,
	}
}
func (s *UserService) CreateUser(ctx context.Context, email string, password string) (*models.User, error) {
	const op = "services.auth.CreateUser"

	log := s.log.With(slog.String("op", op))

	var cost int = 14
	if s.config.Env == "local" {
		cost = 4
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userSaver.CreateUser(ctx, email, string(bytes))
	if err != nil {
		if errors.Is(err, storage.UserExists) {
			log.Info("email is taken")
			return nil, service.ErrEmailExists
		}

		log.Error("couldn't save the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user has been created")

	return user, nil
}
