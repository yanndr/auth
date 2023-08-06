package pg

import (
	"auth/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Open(configuration config.Database) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		configuration.Host, configuration.Port, configuration.UserName, configuration.Password, configuration.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Query("select * from version;")
	if err != nil {
		return nil, fmt.Errorf("db schema not present: %w", err)
	}

	return db, nil
}
