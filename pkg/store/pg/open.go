package pg

import (
	"auth/pkg/model"
	embed "auth/sql/postgresql"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Open(configuration model.DatabaseConfiguration) (*sql.DB, error) {
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
		if _, err := db.ExecContext(context.Background(), embed.Schema); err != nil {
			return nil, err
		}
	}

	return db, nil
}
