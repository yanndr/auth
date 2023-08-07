//go:build pg_test
// +build pg_test

package integrations

import (
	"auth/pkg/config"
	"auth/pkg/stores"
	"auth/pkg/stores/pg"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func openPgDb() (*sql.DB, stores.UserStore, func(), error) {
	database, err := pg.Open(config.Database{
		Host:     "localhost",
		Port:     5433,
		UserName: "auth_user",
		Password: "autPassw@ord",
		DbName:   "auth",
		SslMode:  "disable",
	})
	if err != nil {
		return nil, nil, nil, err
	}
	tx, err := database.Begin()
	if err != nil {
		return nil, nil, nil, err
	}

	tearDown := func() {
		tx.Rollback()
	}

	return database, stores.NewPgUserStore(pg.New(tx)), tearDown, nil
}

func Test_pg_Server_Create(t *testing.T) {
	teardown := setup(t, openPgDb)
	defer teardown(t)
	testServerCreate(t)
}

func Test_pg_Server_Create_SameUsername(t *testing.T) {
	teardown := setup(t, openPgDb)
	defer teardown(t)
	testServerCreateSameUser(t)
}

func Test_pg_Server_Auth_Success(t *testing.T) {
	teardown := setup(t, openPgDb)
	defer teardown(t)
	testServerAuthSuccess(t)
}

func Test_pg_Server_Auth_Failure(t *testing.T) {
	teardown := setup(t, openPgDb)
	defer teardown(t)
	testServerAuthFailure(t)
}
