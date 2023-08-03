package store

import (
	"auth/pkg/model"
	"auth/pkg/store/pg"
	"context"
	"database/sql"
)

type UserStore interface {
	Create(ctx context.Context, user model.User) error
	Get(ctx context.Context, username string) (*model.User, error)
}

type userStore struct {
	querier pg.Querier
}

func NewUserStore(q pg.Querier) UserStore {
	return &userStore{querier: q}
}

func (s *userStore) Create(ctx context.Context, user model.User) error {
	_, err := s.querier.CreateUser(ctx, pg.CreateUserParams{Username: user.Username, Password: user.Password})
	if err != nil {
		return err
	}

	return nil
}

func (s *userStore) Get(ctx context.Context, username string) (*model.User, error) {

	u, err := s.querier.GetUser(ctx, username)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}

	}

	return &model.User{Username: username, Password: u.Password}, nil
}
