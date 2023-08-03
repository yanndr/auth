package services

import (
	"auth/internal/model"
	"auth/internal/store"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	Create(ctx context.Context, user model.User) error
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	userStore store.UserStore
}

func NewUserService(userStore store.UserStore) UserService {
	return &userService{
		userStore: userStore,
	}
}

func (s *userService) Create(ctx context.Context, user model.User) error {
	u, err := s.userStore.Get(ctx, user.Username)
	if err != nil {
		return err
	}
	if u != nil {
		return fmt.Errorf("user already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = s.userStore.Create(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (string, error) {
	u, err := s.userStore.Get(ctx, username)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", fmt.Errorf("user %s doesn't exist", username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := generateJWT(u)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateJWT(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["sub"] = user.Username
	claims["iss"] = "authService"
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString([]byte("captainjacksparrowsayshi"))

	if err != nil {
		_ = fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
