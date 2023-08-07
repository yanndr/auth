package stores

import (
	"auth/pkg/config"
	"auth/pkg/stores/pg"
	"database/sql"
	"fmt"
	"testing"
)

func openPgDb() (*sql.DB, error) {
	database, err := pg.Open(config.Database{
		Host:     "localhost",
		Port:     5433,
		UserName: "auth_user",
		Password: "autPassw@ord",
		DbName:   "auth",
		SslMode:  "disable",
	})
	if err != nil {
		return nil, fmt.Errorf("an error %v was not expected when opening a test database connection", err)
	}
	return database, nil
}

func TestPgUserStore_Create(t *testing.T) {
	testCreateUser(t, openPgDb)
}

func TestPgUserStore_Get(t *testing.T) {
	testGetUser(t, openPgDb)
}
