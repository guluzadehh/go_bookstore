package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/guluzadehh/go_bookstore/services/user/internal/domain/models"
	"github.com/guluzadehh/go_bookstore/services/user/internal/storage"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateUser(ctx context.Context, email string, password string) (*models.User, error) {
	const op = "storage.postgresql.CreateUser"

	var user models.User

	const query = `
		INSERT INTO users(email, password)
		VALUES ($1, $2)
		RETURNING id, email, password, created_at, updated_at, is_active
	`

	var p string

	err := s.db.QueryRowContext(ctx, query, email, password).Scan(
		&user.Id,
		&user.Email,
		&p,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code.Name() == "unique_violation" {
			return nil, storage.UserExists
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Password = []byte(p)

	return &user, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgresql.UserByEmail"

	var user models.User

	const query = `
		SELECT users.id, users.email, users.password, users.created_at, users.updated_at, users.is_active
		FROM users
		WHERE users.email = $1;
	`

	if err := s.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.UserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
