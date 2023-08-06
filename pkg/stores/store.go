package stores

import (
	"auth/pkg/models"
	"context"
)

type UserStore interface {
	Create(ctx context.Context, user models.User) error
	Get(ctx context.Context, username string) (*models.User, error)
}
