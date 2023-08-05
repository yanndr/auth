package services

import (
	autherrors "auth/pkg/errors"
	"auth/pkg/jwt"
	"auth/pkg/store"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Authentication interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type JwtAuthenticationService struct {
	UserStore    store.UserStore
	JwtGenerator jwt.Generator
}

func (as *JwtAuthenticationService) Authenticate(ctx context.Context, username, password string) (string, error) {
	u, err := as.UserStore.Get(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error getting the user %s: %w", username, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", autherrors.AutenticationFailErr
		}
		return "", fmt.Errorf("error comparing password: %w", err)
	}
	token, err := as.JwtGenerator.GenerateJWT(u)
	if err != nil {
		return "", fmt.Errorf("error generating the token: %w", err)
	}
	return token, nil
}