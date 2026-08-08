package database

import "database/sql"

var scheme = `
	CREATE TABLE IF NOT EXISTS movie(
	id 			INTEGER PRIMARY KEY AUTOINCREMENT, 
	title       TEXT NOT NULL UNIQUE,
	releaseYear INTEGER NOT NULL,
	duration    INTEGER  NOT NULL,
	created_at 	DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at 	DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`

func InitializeMovieTable(db *sql.DB) error {

	_, err := db.Exec(scheme)
	return err
}
