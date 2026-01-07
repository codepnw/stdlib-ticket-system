package userrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/user"
	"github.com/lib/pq"
)

type UserRepository interface {
	CreateUser(ctx context.Context, input user.User) (user.User, error)
	FindUsername(ctx context.Context, username string) (user.User, error)
	
	// Auth
	SaveRefreshToken(ctx context.Context, input user.Auth) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, input user.User) (user.User, error) {
	query := `
		INSERT INTO users (username, hash_password)
		VALUES ($1, $2) RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, input.Username, input.HashPassword).Scan(
		&input.ID,
	)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code == pq.ErrorCode("23505") {
			return user.User{}, errs.ErrUsernameAlreadyExists
		}
		return user.User{}, err
	}
	return input, nil
}

func (r *userRepository) FindUsername(ctx context.Context, username string) (user.User, error) {
	query := `
		SELECT id, username, hash_password
		FROM users WHERE username = $1 LIMIT 1
	`
	var u user.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID,
		&u.Username,
		&u.HashPassword,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, errs.ErrInvalidCredentials
		}
		return user.User{}, nil
	}
	return u, nil
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, input user.Auth) error {
	query := `
		INSERT INTO auth (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			refresh_token = EXCLUDED.refresh_token, expires_at = EXCLUDED.expires_at, revoked = FALSE
	`
	_, err := r.db.ExecContext(ctx, query, input.UserID, input.RefreshToken, input.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}
