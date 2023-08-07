package services

import (
	autherrors "auth/pkg/errors"
	"auth/pkg/models"
	"auth/pkg/stores"
	"auth/pkg/validators"
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user models.User) error
}

type userService struct {
	userStore stores.UserStore
	validator validators.Validator
	hashCost  int
	logger    *zap.Logger
}

func NewUserService(userStore stores.UserStore, validator validators.Validator, hashCost int) UserService {
	return &userService{
		userStore: userStore,
		validator: validator,
		hashCost:  hashCost,
		logger:    zap.L().Named("UserService"),
	}
}

func (s *userService) Create(ctx context.Context, userRequest models.User) error {
	err := s.validator.Validate(userRequest)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	u, err := s.userStore.Get(ctx, userRequest.Username)
	if err != nil {
		return fmt.Errorf("error getting user from store: %w", err)
	}
	if u != nil {
		return autherrors.UsernameAlreadyExistErr{Name: userRequest.Username}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), s.hashCost)
	if err != nil {
		return fmt.Errorf("error during password hashing: %w", err)
	}
	user := models.User{
		Username: userRequest.Username,
		Password: string(hashedPassword),
	}

	err = s.userStore.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating the user: %w", err)
	}

	return nil
}
