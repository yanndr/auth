package store

import (
	"auth/internal/model"
	"auth/internal/store/pg"
	"context"
	"database/sql"
)

type UserStore interface {
	Create(user model.User) error
	Get(username string) (*model.User, error)
}

type userStore struct {
	querier pg.Querier
}

func NewUserStore(q pg.Querier) UserStore {
	return &userStore{querier: q}
}

func (s *userStore) Create(user model.User) error {
	_, err := s.querier.CreateUser(context.Background(), pg.CreateUserParams{Username: user.Username, Password: user.Password})
	if err != nil {
		return err
	}

	return nil
}

func (s *userStore) Get(username string) (*model.User, error) {

	u, err := s.querier.GetUser(context.Background(), username)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}

	}

	return &model.User{Username: username, Password: u.Password}, nil
}
