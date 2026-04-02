package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Optional: connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
