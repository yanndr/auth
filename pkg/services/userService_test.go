package services

import (
	autherrors "auth/pkg/errors"
	"auth/pkg/models"
	"auth/pkg/tests"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	mockUserStore    *tests.MockUserStore
	mockValidator    *tests.MockValidator
	mockJwtGenerator *tests.MockTokenGenerator
)

func setupTest(t testing.TB) func(t testing.TB) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserStore = tests.NewMockUserStore(ctrl)
	mockValidator = tests.NewMockValidator(ctrl)
	mockJwtGenerator = tests.NewMockTokenGenerator(ctrl)

	return func(t testing.TB) {
	}
}

func Test_userService_Create_no_error(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	user := models.User{Username: "test", Password: "test"}

	mockValidator.EXPECT().Validate(user).Return(nil).Times(1)
	mockUserStore.EXPECT().Get(ctx, user.Username).Times(1)
	mockUserStore.EXPECT().Create(ctx, gomock.Any()).Times(1)

	s := &userService{
		userStore: mockUserStore,
		validator: mockValidator,
	}

	err := s.Create(ctx, user)
	require.NoError(t, err)
}

func Test_userService_Create_validator_error(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	user := models.User{Username: "test", Password: "test"}

	errorMsg := "something is not valid"
	mockValidator.EXPECT().Validate(user).Return(fmt.Errorf(errorMsg)).Times(1)
	mockUserStore.EXPECT().Get(ctx, user.Username).Times(0)
	mockUserStore.EXPECT().Create(ctx, gomock.Any()).Times(0)

	s := &userService{
		userStore: mockUserStore,
		validator: mockValidator,
	}

	err := s.Create(ctx, user)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("validation error: %s", errorMsg))
}

func Test_userService_Create_store_get_error(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	user := models.User{Username: "test", Password: "test"}

	errorMsg := "something went wrong"
	mockValidator.EXPECT().Validate(user).Return(nil).Times(1)
	mockUserStore.EXPECT().Get(ctx, user.Username).Return(nil, fmt.Errorf(errorMsg)).Times(1)
	mockUserStore.EXPECT().Create(ctx, gomock.Any()).Times(0)

	s := &userService{
		userStore: mockUserStore,
		validator: mockValidator,
	}

	err := s.Create(ctx, user)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("error getting user from store: %s", errorMsg))
}

func Test_userService_Create_store_get_existing_user(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	user := models.User{Username: "test", Password: "test"}
	existingUser := models.User{Username: "test", Password: "jkljkljkljkljkl"}

	mockValidator.EXPECT().Validate(user).Return(nil).Times(1)
	mockUserStore.EXPECT().Get(ctx, user.Username).Return(&existingUser, nil).Times(1)
	mockUserStore.EXPECT().Create(ctx, gomock.Any()).Times(0)

	s := &userService{
		userStore: mockUserStore,
		validator: mockValidator,
	}

	err := s.Create(ctx, user)
	require.Error(t, err)
	require.EqualError(t, err, autherrors.UsernameAlreadyExistErr{Name: user.Username}.Error())
}

func Test_userService_Create_store_create_error(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	user := models.User{Username: "test", Password: "test"}
	errorMsg := "something went wrong"

	mockValidator.EXPECT().Validate(user).Return(nil).Times(1)
	mockUserStore.EXPECT().Get(ctx, user.Username).Return(nil, nil).Times(1)
	mockUserStore.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf(errorMsg)).Times(1)

	s := &userService{
		userStore: mockUserStore,
		validator: mockValidator,
	}

	err := s.Create(ctx, user)
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("error creating the user: %s", errorMsg))
}
