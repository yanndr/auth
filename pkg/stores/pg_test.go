//go:build pg_test

package stores

import (
	"auth/pkg/config"
	"auth/pkg/stores/pg"
	"testing"
)

func setupPg(t testing.TB) func(t testing.TB) {
	database, err := pg.Open(config.Database{
		Host:     "localhost",
		Port:     5433,
		UserName: "auth_user",
		Password: "autPassw@ord",
		DbName:   "auth",
		SslMode:  "disable",
	})
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

func TestPgUserStore_Create(t *testing.T) {
	testCreateUser(t, setupPg)
}

func TestPgUserStore_Get(t *testing.T) {
	testGetUser(t, setupPg)
}
