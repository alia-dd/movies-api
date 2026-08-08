package database

import (
	"database/sql"
)

const schema = `
	CREATE TABLE IF NOT EXISTS ACTORS(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	birthdate TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

func NewTable(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
