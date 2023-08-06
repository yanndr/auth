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
	"database/sql"
	"fmt"
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

func inMemoryUserStore() (*sql.DB, stores.UserStore, error) {
	database, err := sqlite.OpenInMemory()
	if err != nil {
		return nil, nil, fmt.Errorf("an error %v was not expected when opening a stub database connection", err)
	}
	return database, stores.NewSqliteUserStore(sqlite.New(database)), nil
}

func setup(t testing.TB, storeFn func() (*sql.DB, stores.UserStore, error)) func(t testing.TB) {

	var err error
	var database *sql.DB
	database, userStore, err = storeFn()
	if err != nil {
		t.Fatalf("an error %v was not expected when opening a stub database connection", err)
	}

	userValidator := validators.UserValidator{
		StructValidator:   validator.New(),
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
	teardown := setup(t, inMemoryUserStore)
	defer teardown(t)
	testServerCreate(t)
}

func Test_Server_Create_SameUsername(t *testing.T) {
	teardown := setup(t, inMemoryUserStore)
	defer teardown(t)
	testServerCreateSameUser(t)
}

func Test_Server_Auth_Success(t *testing.T) {
	teardown := setup(t, inMemoryUserStore)
	defer teardown(t)
	testServerAuthSuccess(t)
}

func Test_Server_Auth_Failure(t *testing.T) {
	teardown := setup(t, inMemoryUserStore)
	defer teardown(t)
	testServerAuthFailure(t)
}

func testServerCreate(t *testing.T) {
	username := "test"
	password := "password"

	createUser(t, username, password)

	u, err := userStore.Get(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, username, u.Username)
	require.NotEmpty(t, u.Password)
	require.NotEqual(t, password, u.Password)
}

func testServerAuthSuccess(t *testing.T) {
	username := "test3"
	password := "password"
	createUser(t, username, password)

	response, err := grpcServer.Authenticate(context.Background(), &pb.AuthenticateRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.NotEmpty(t, response.Token)
}

func testServerAuthFailure(t *testing.T) {
	username := "test4"
	password := "password"
	createUser(t, username, password)

	response, err := grpcServer.Authenticate(context.Background(), &pb.AuthenticateRequest{
		Username: username,
		Password: "notright",
	})
	require.Error(t, err)
	require.Empty(t, response)
}

func testServerCreateSameUser(t *testing.T) {
	username := "test2"
	password := "password"
	createUser(t, username, password)

	response, err := grpcServer.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})
	require.Error(t, err)
	require.Empty(t, response)
	require.EqualError(t, err, "rpc error: code = AlreadyExists desc = username test2 already exists")
}

func createUser(t *testing.T, username string, password string) {
	response, err := grpcServer.CreateUser(context.Background(), &pb.CreateUserRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	_, err = userStore.Get(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, true, response.Success)
}
