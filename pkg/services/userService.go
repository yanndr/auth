package services

import (
	errors2 "auth/pkg/errors"
	"auth/pkg/models"
	"auth/pkg/store"
	"auth/pkg/validators"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user models.User) error
}

type userService struct {
	userStore store.UserStore
	validator validators.Validator
}

func NewUserService(userStore store.UserStore, validator validators.Validator) UserService {
	return &userService{
		userStore: userStore,
		validator: validator,
	}
}

func (s *userService) Create(ctx context.Context, user models.User) error {
	err := s.validator.Validate(user)
	if err != nil {
		return err
	}

	u, err := s.userStore.Get(ctx, user.Username)
	if err != nil {
		return err
	}
	if u != nil {
		return errors2.UsernameAlreadyExistErr{Name: user.Username}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	err = s.userStore.Create(ctx, user)

	if err != nil {
		return err
	}

	return nil
}
