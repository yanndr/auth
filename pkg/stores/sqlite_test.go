package stores

import (
	"auth/pkg/models"
	"auth/pkg/stores/pg"
	"auth/pkg/stores/sqlite"
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

var (
	userStore UserStore
)

func setup(t testing.TB, openDbFn func() (*sql.DB, error)) func(t testing.TB) {
	database, err := openDbFn()
	if err != nil {
		t.Fatalf("an error %v was not expected when opening a test database connection", err)
	}
	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("an error %v was not expected when opening a stub database transaction", err)
	}
	userStore = NewPgUserStore(pg.New(tx))

	return func(t testing.TB) {
		tx.Rollback()
		database.Close()
	}
}

func openSqliteDb() (*sql.DB, error) {
	database, err := sqlite.OpenInMemory()
	if err != nil {
		return nil, fmt.Errorf("an error %v was not expected when opening a test database connection", err)
	}
	return database, nil
}

func TestSqliteUserStore_Create(t *testing.T) {
	testCreateUser(t, openSqliteDb)
}

func TestSqliteUserStore_Get(t *testing.T) {
	testGetUser(t, openSqliteDb)
}

func testCreateUser(t *testing.T, openDbFn func() (*sql.DB, error)) {
	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{"Valid user", models.User{"test", "sdfafdfasdfds"}, false},
		{"No username", models.User{"", "sdfafdfasdfds"}, true},
		{"No password", models.User{"tesr", ""}, true},
		{"Nothing", models.User{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup(t, openDbFn)
			defer teardown(t)
			s := userStore
			if err := s.Create(context.Background(), tt.user); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func testGetUser(t *testing.T, openDbFn func() (*sql.DB, error)) {
	user := models.User{"test", "fsdjak"}
	tests := []struct {
		name         string
		existingUser models.User
		username     string
		want         *models.User
		wantErr      bool
	}{
		{"get existing user", user, "test", &user, false},
		{"get  non existing user", user, "test2", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup(t, openDbFn)
			defer teardown(t)
			s := userStore
			_ = s.Create(context.Background(), tt.existingUser)
			got, err := s.Get(context.Background(), tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
