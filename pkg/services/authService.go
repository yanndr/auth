package services

import (
	autherrors "auth/pkg/errors"
	"auth/pkg/jwt"
	"auth/pkg/stores"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	//Authenticate a user from a username and password
	Authenticate(ctx context.Context, username, password string) (string, error)
}

// JwtAuthService is an implementation of AuthService that returns a JWT.
type JwtAuthService struct {
	UserStore    stores.UserStore
	JwtGenerator jwt.TokenGenerator
	logger       *zap.Logger
}

// NewJwtAuthService creates a new instance of an AuthService using JWT
func NewJwtAuthService(userStore stores.UserStore, jwtGenerator jwt.TokenGenerator) AuthService {
	return &JwtAuthService{
		UserStore:    userStore,
		JwtGenerator: jwtGenerator,
		logger:       zap.L().Named("AuthService"),
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
		as.logger.Error("failed to compare passwords", zap.Error(err))
		return "", fmt.Errorf("error comparing password: %w", err)
	}

	token, err := as.JwtGenerator.Generate(*u)
	if err != nil {
		return "", fmt.Errorf("error generating the token: %w", err)
	}
	return token, nil
}
