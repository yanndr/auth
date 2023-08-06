package stores

import (
	"auth/pkg/models"
	"auth/pkg/stores/pg"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PgUserStore struct {
	querier pg.Querier
}

func NewPgUserStore(q pg.Querier) UserStore {
	return &PgUserStore{querier: q}
}

func (s *PgUserStore) Create(ctx context.Context, user models.User) error {
	_, err := s.querier.CreateUser(ctx, pg.CreateUserParams{Username: user.Username, PasswordHash: user.Password})
	if err != nil {
		return fmt.Errorf("error creating the user %s: %w", user.Username, err)
	}

	return nil
}

func (s *PgUserStore) Get(ctx context.Context, username string) (*models.User, error) {

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
