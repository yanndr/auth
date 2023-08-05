package store

import (
	"auth/pkg/models"
	"auth/pkg/store/pg"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserStore interface {
	Create(ctx context.Context, user models.User) error
	Get(ctx context.Context, username string) (*models.User, error)
}

type userStore struct {
	querier pg.Querier
}

func NewUserStore(q pg.Querier) UserStore {
	return &userStore{querier: q}
}

func (s *userStore) Create(ctx context.Context, user models.User) error {
	_, err := s.querier.CreateUser(ctx, pg.CreateUserParams{Username: user.Username, PasswordHash: user.PasswordHash})
	if err != nil {
		return fmt.Errorf("error cerating the user %s: %w", user.Username, err)
	}

	return nil
}

func (s *userStore) Get(ctx context.Context, username string) (*models.User, error) {

	u, err := s.querier.GetUser(ctx, username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error getting the user %s: %w", username, err)
		} else {
			return nil, nil
		}

	}

	return &models.User{Username: username, PasswordHash: u.PasswordHash}, nil
}
