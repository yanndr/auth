package services

import (
	autherrors "auth/pkg/errors"
	"auth/pkg/models"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_authService_Authenticate_no_error(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	ctx := context.Background()
	password := "test"
	username := "user"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{Username: username, Password: string(hash)}

	mockUserStore.EXPECT().Get(ctx, username).Return(&user, nil).Times(1)
	mockJwtGenerator.EXPECT().Generate(user).Return("sdjklfjasdkl.jfsda.fasdf", nil).Times(1)

	s := &JwtAuthService{
		UserStore:    mockUserStore,
		JwtGenerator: mockJwtGenerator,
	}

	token, err := s.Authenticate(ctx, username, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func Test_authService_Authenticate_store_get_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	//Prepare
	ctx := context.Background()
	password := "test"
	username := "user"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{Username: username, Password: string(hash)}
	errorMsg := "something went wrong"

	mockUserStore.EXPECT().Get(ctx, username).Return(nil, fmt.Errorf(errorMsg)).Times(1)
	mockJwtGenerator.EXPECT().Generate(&user).Return("sdjklfjasdkl.jfsda.fasdf", nil).Times(0)

	s := &JwtAuthService{
		UserStore:    mockUserStore,
		JwtGenerator: mockJwtGenerator,
	}

	//Act
	token, err := s.Authenticate(ctx, username, password)

	//Verify
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("error getting user %s from store: %s", username, errorMsg))
	require.Empty(t, token)
}

func Test_authService_Authenticate_store_no_user(t *testing.T) {

	teardownTest := setupTest(t)
	defer teardownTest(t)

	//Prepare
	ctx := context.Background()
	password := "test"
	username := "user"

	mockUserStore.EXPECT().Get(ctx, username).Return(nil, nil).Times(1)
	mockJwtGenerator.EXPECT().Generate(gomock.Any()).Times(0)

	s := &JwtAuthService{
		UserStore:    mockUserStore,
		JwtGenerator: mockJwtGenerator,
	}

	//Act
	token, err := s.Authenticate(ctx, username, password)

	//Verify
	require.Error(t, err)
	require.EqualError(t, err, autherrors.AuthenticationFailErr(username).Error())
	require.Empty(t, token)
}

func Test_authService_Authenticate_password_mismatch(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	//Prepare
	ctx := context.Background()
	password := "test"
	username := "user"
	hash, _ := bcrypt.GenerateFromPassword([]byte("other"), 10)
	user := models.User{Username: username, Password: string(hash)}

	mockUserStore.EXPECT().Get(ctx, username).Return(&user, nil).Times(1)
	mockJwtGenerator.EXPECT().Generate(&user).Times(0)

	s := &JwtAuthService{
		UserStore:    mockUserStore,
		JwtGenerator: mockJwtGenerator,
	}

	//Act
	token, err := s.Authenticate(ctx, username, password)

	//Verify
	require.Error(t, err)
	require.EqualError(t, err, autherrors.AuthenticationFailErr(username).Error())
	require.Empty(t, token)
}

func Test_authService_Authenticate_bcrypt_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	//Prepare
	ctx := context.Background()
	password := "test"
	username := "user"

	user := models.User{Username: username, Password: ""}

	mockUserStore.EXPECT().Get(ctx, username).Return(&user, nil).Times(1)
	mockJwtGenerator.EXPECT().Generate(&user).Times(0)

	s := &JwtAuthService{
		UserStore:    mockUserStore,
		JwtGenerator: mockJwtGenerator,
	}

	//Act
	token, err := s.Authenticate(ctx, username, password)

	//Verify
	require.Error(t, err)
	require.EqualError(t, err, "error comparing password: crypto/bcrypt: hashedSecret too short to be a bcrypted password")
	require.Empty(t, token)

}

func Test_authService_Authenticate_generator_error(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	//Prepare
	ctx := context.Background()
	password := "test"
	username := "user"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{Username: username, Password: string(hash)}
	errorMsg := "can't generate the token"

	mockUserStore.EXPECT().Get(ctx, username).Return(&user, nil).Times(1)
	mockJwtGenerator.EXPECT().Generate(user).Return("", fmt.Errorf(errorMsg)).Times(1)

	s := &JwtAuthService{
		UserStore:    mockUserStore,
		JwtGenerator: mockJwtGenerator,
	}

	//Act
	token, err := s.Authenticate(ctx, username, password)

	//Verify
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("error generating the token: %s", errorMsg))
	require.Empty(t, token)
}
