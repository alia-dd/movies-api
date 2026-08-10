package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(db string) (*sql.DB, error) {
	dbName := db + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON"
	dataBase, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return dataBase, err
	}
	dataBase.SetMaxOpenConns(1) // SQLite supports only one writer at a time
	dataBase.SetMaxIdleConns(1)
	dataBase.SetConnMaxLifetime(time.Hour)
	return dataBase, nil

}
