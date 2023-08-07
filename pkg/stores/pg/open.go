package pg

import (
	"auth/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Open(configuration config.Database) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s sslkey=%s sslcert=%s",
		configuration.Host,
		configuration.Port,
		configuration.UserName,
		configuration.Password,
		configuration.DbName,
		configuration.SslMode,
		configuration.RootCert,
		configuration.SslKey,
		configuration.SslCert,
	)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("connection to database not alive: %w", err)
	}

	_, err = db.Query("select * from version;")
	if err != nil {
		return nil, fmt.Errorf("db schema not present: %w", err)
	}

	return db, nil
}
