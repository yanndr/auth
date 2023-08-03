package pg

import (
	embed "auth/sql/postgresql"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Open(host string, port int, user, password, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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
