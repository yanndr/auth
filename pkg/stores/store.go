package stores

import (
	"auth/pkg/models"
	"context"
)

type UserStore interface {
	//Create a user from models.User and store it.
	Create(ctx context.Context, user models.User) error
	//Get a user with the username from the store.
	Get(ctx context.Context, username string) (*models.User, error)
}
