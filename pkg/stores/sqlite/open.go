package sqlite

import (
	"auth/sql/sqlite"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func Open(dbPath string) (*sql.DB, error) {
	newDb := false
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		newDb = true
	}
	database, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbPath))
	if err != nil {
		return nil, err
	}

	if newDb {
		if _, err := database.ExecContext(context.Background(), sqlite.Schema); err != nil {
			return nil, err
		}
	}

	return database, nil
}

func OpenInMemory() (*sql.DB, error) {

	database, err := sql.Open("sqlite3", fmt.Sprintf("file::memory:?_foreign_keys=on"))
	if err != nil {
		return nil, err
	}

	if _, err := database.ExecContext(context.Background(), sqlite.Schema); err != nil {
		return nil, err
	}

	return database, nil
}
