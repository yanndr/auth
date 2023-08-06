package server

import (
	"auth/pkg/errors"
	"auth/pkg/models"
	"auth/pkg/pb"
	"auth/pkg/tests"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	mockAuthentication *tests.MockAuthService
	mockUserService    *tests.MockUserService
)

func setupTest(t testing.TB) func(t testing.TB) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthentication = tests.NewMockAuthService(ctrl)
	mockUserService = tests.NewMockUserService(ctrl)

	return func(t testing.TB) {
	}
}

func TestAuthServer_CreateUser_no_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	username := "test"
	password := "password"
	user := models.User{
		Username: username,
		Password: password,
	}
	mockUserService.EXPECT().Create(ctx, user).Return(nil).Times(1)
	server := NewServer(mockUserService, mockAuthentication)

	response, err := server.CreateUser(ctx, &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})

	require.NoError(t, err)
	require.Equal(t, response.Success, true)
}

func TestAuthServer_CreateUser_create_username_exists(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	username := "test"
	password := "password"
	user := models.User{
		Username: username,
		Password: password,
	}
	mockUserService.EXPECT().Create(ctx, user).Return(errors.UsernameAlreadyExistErr{Name: username}).Times(1)
	server := NewServer(mockUserService, mockAuthentication)

	response, err := server.CreateUser(ctx, &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})

	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = AlreadyExists desc = username test already exists")
	require.Empty(t, response)
}

func TestAuthServer_CreateUser_create_username_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	username := "test"
	password := "password"
	user := models.User{
		Username: username,
		Password: password,
	}
	mockUserService.EXPECT().Create(ctx, user).Return(fmt.Errorf("unexpected")).Times(1)
	server := NewServer(mockUserService, mockAuthentication)

	response, err := server.CreateUser(ctx, &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})

	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = Unknown desc = unexpected")
	require.Empty(t, response)
}

func TestAuthServer_Authenticate_no_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	username := "test"
	password := "password"

	mockAuthentication.EXPECT().Authenticate(ctx, username, password).Return("token", nil).Times(1)
	server := NewServer(mockUserService, mockAuthentication)

	response, err := server.Authenticate(ctx, &pb.AuthenticateRequest{
		Username: username,
		Password: password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, response.Token)
	require.Equal(t, "token", response.Token)
}

func TestAuthServer_Authenticate_authService_authFailed(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	username := "test"
	password := "password"

	mockAuthentication.EXPECT().Authenticate(ctx, username, password).Return("", errors.AuthenticationFailErr(username)).Times(1)
	server := NewServer(mockUserService, mockAuthentication)

	response, err := server.Authenticate(ctx, &pb.AuthenticateRequest{
		Username: username,
		Password: password,
	})

	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = authentication failed")
	require.Empty(t, response)
}

func TestAuthServer_Authenticate_authService_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	username := "test"
	password := "password"

	mockAuthentication.EXPECT().Authenticate(ctx, username, password).Return("", fmt.Errorf("unexpected")).Times(1)
	server := NewServer(mockUserService, mockAuthentication)

	response, err := server.Authenticate(ctx, &pb.AuthenticateRequest{
		Username: username,
		Password: password,
	})

	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = Unknown desc = unexpected")
	require.Empty(t, response)
}
