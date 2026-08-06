package database

import (
	"database/sql"
)

func NewTable(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS ACTORS(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	birthdate TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	return err
}
