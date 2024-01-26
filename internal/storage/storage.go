package storage

import (
	"database/sql"
)

func LoadSQLDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}

func CLoseSQLDB(db *sql.DB) error {
	return db.Close()
}
