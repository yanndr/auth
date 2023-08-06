package integrations

import (
	"auth/pkg/config"
	"auth/pkg/jwt"
	"auth/pkg/pb"
	"auth/pkg/server"
	"auth/pkg/services"
	"auth/pkg/stores"
	"auth/pkg/stores/sqlite"
	"auth/pkg/validators"
	"context"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	grpcServer  *server.AuthServer
	userService services.UserService
	authService services.AuthService
	userStore   stores.UserStore
)

func setup(t testing.TB) func(t testing.TB) {

	database, err := sqlite.OpenInMemory()
	if err != nil {
		t.Fatalf("an error %v was not expected when opening a stub database connection", err)
	}

	userStore = stores.NewSqliteUserStore(sqlite.New(database))
	userValidator := validators.UserValidator{
		Validator:         validator.New(),
		PasswordValidator: validators.NewPasswordValidator(config.Password{}),
	}
	userService = services.NewUserService(userStore, userValidator, 10)
	jwtGenerator := jwt.NewGenerator(config.Token{
		SigningMethod: "HS256",
		SignedKey:     "sdfsadfa",
		Audience:      "audience",
		Issuer:        "issuer",
		ExpDuration:   10,
	})
	authService = services.NewJwtAuthService(userStore, jwtGenerator)

	grpcServer = server.NewServer(userService, authService)

	return func(t testing.TB) {
		database.Close()
	}
}

func Test_Server_Create(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)
	username := "test"
	password := "password"

	response, err := grpcServer.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, true, response.Success)

	u, err := userStore.Get(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, username, u.Username)
	require.NotEmpty(t, u.Password)
	require.NotEqual(t, password, u.Password)
}

func Test_userService_Create_SameUsername(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)
	username := "test"
	password := "password"
	_, err := grpcServer.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	_, err = userStore.Get(context.Background(), username)
	require.NoError(t, err)

	response, err := grpcServer.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})
	require.Error(t, err)
	require.Empty(t, response)
	require.EqualError(t, err, "rpc error: code = AlreadyExists desc = username test already exists")
}
