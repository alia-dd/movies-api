package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
)

func InitializeDB(dbName string) (*sql.DB, error) {
	// this will try to open the database, If the database doesn’t exist, it will automatically be created.
	dbConnStr := dbName + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON"
	// db, dbErr := sql.Open("sqlite", dbConnStr)
	db, dbErr := sql.Open("sqlite3", dbConnStr)
	if dbErr != nil {
		return nil, fmt.Errorf("Failed to open database: %w", dbErr)
	}

	db.SetMaxOpenConns(1) // SQLite supports only one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	if errPing := db.Ping(); errPing != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", errPing)
	}
	log.Printf("Successfully connected to %s database", dbName)

	return db, nil
}

// func OpenDB(db string) (*sql.DB, error) {
// 	dbName := db + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON"
// 	dataBase, err := sql.Open("sqlite3", dbName)
// 	if err != nil {
// 		return dataBase, err
// 	}
// 	dataBase.SetMaxOpenConns(1) // SQLite supports only one writer at a time
// 	dataBase.SetMaxIdleConns(1)
// 	dataBase.SetConnMaxLifetime(time.Hour)
// 	return dataBase, nil

// }
