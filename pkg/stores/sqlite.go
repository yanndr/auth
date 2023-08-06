package stores

import (
	"auth/pkg/models"
	"auth/pkg/stores/sqlite"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SqliteUserStore struct {
	querier sqlite.Querier
}

func NewSqliteUserStore(q sqlite.Querier) UserStore {
	return &SqliteUserStore{querier: q}
}

func (s *SqliteUserStore) Create(ctx context.Context, user models.User) error {
	_, err := s.querier.CreateUser(ctx, sqlite.CreateUserParams{Username: user.Username, PasswordHash: user.Password})
	if err != nil {
		return fmt.Errorf("error creating the user %s: %w", user.Username, err)
	}

	return nil
}

func (s *SqliteUserStore) Get(ctx context.Context, username string) (*models.User, error) {

	u, err := s.querier.GetUser(ctx, username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error getting the user %s: %w", username, err)
		} else {
			return nil, nil
		}
	}

	return &models.User{Username: u.Username, Password: u.PasswordHash}, nil
}
