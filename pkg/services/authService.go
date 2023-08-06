package services

import (
	autherrors "auth/pkg/errors"
	"auth/pkg/jwt"
	"auth/pkg/stores"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type JwtAuthService struct {
	UserStore    stores.UserStore
	JwtGenerator jwt.JwtGenerator
}

func NewJwtAuthService(userStore stores.UserStore, jwtGenerator jwt.JwtGenerator) AuthService {
	return &JwtAuthService{
		UserStore:    userStore,
		JwtGenerator: jwtGenerator,
	}
}

func (as *JwtAuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	u, err := as.UserStore.Get(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error getting user %s from store: %w", username, err)
	}
	if u == nil {
		return "", autherrors.AuthenticationFailErr(username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", autherrors.AuthenticationFailErr(u.Username)
		}
		return "", fmt.Errorf("error comparing password: %w", err)
	}

	token, err := as.JwtGenerator.Generate(*u)
	if err != nil {
		return "", fmt.Errorf("error generating the token: %w", err)
	}
	return token, nil
}