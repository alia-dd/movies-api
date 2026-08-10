package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movies-api/internal/models"
	"os"
)

var scheme = `
	CREATE TABLE IF NOT EXISTS movie(
	id 					INTEGER PRIMARY KEY AUTOINCREMENT, 
	title       		TEXT NOT NULL UNIQUE,
	releaseYear 		INTEGER NOT NULL,
	duration    		INTEGER,
	Overview 			TEXT,
	OriginalLanguage 	TEXT,
	created_at 			DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at 			DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS ACTORS(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	birthdate TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS movie_genre(
	genreId      		INTEGER NOT NULL,
	movieId      		INTEGER NOT NULL,
	PRIMARY KEY (movieId, genreId),
    FOREIGN KEY (movieId) REFERENCES movie (id) ,
	FOREIGN KEY (genreId) REFERENCES genre (id) 
	);
	CREATE TABLE IF NOT EXISTS movie_actor(
	actorId      		INTEGER NOT NULL,
	movieId      		INTEGER NOT NULL,
	PRIMARY KEY (movieId, actorId),
    FOREIGN KEY (movieId) REFERENCES movie (id) ,
	FOREIGN KEY (actorId) REFERENCES ACTORS (id) 
	);

`

type MovieSeed struct {
	Results []models.Movies `json:"results"`
}

func InitializeMovieTable(db *sql.DB) error {
	if tableExist(db) {
		return nil
	}
	fmt.Println("false")

	if _, err := db.Exec(scheme); err != nil {
		return fmt.Errorf("failed to create movie table: %w", err)
	}

	if err := seedData(db); err != nil {
		return fmt.Errorf("failed to seed movie data: %w", err)
	}

	return nil
}
func tableExist(db *sql.DB) bool {
	var count int
	query := `SELECT count(*) FROM sqlite_master WHERE name ='movie' and type='table'`
	err := db.QueryRow(query).Scan(&count)
	if err != nil || count <= 0 {
		return false
	}
	return true
}

func NewTable(db *sql.DB) error {
	_, err := db.Exec(scheme)
	return err
}

func seedData(db *sql.DB) error {
	var movies []models.Movies
	query := ` INSERT INTO movie (Title, ReleaseYear, Duration, Overview, OriginalLanguage, GenreId ,ActorId)VALUES (?, ?, ?, ?, ?, ?, ?)`
	mgQuery := `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
	maQuery := `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`

	body, fileErr := os.ReadFile("data/tmdb_movies.json")
	if fileErr != nil {
		return fileErr
	}
	if err := json.Unmarshal(body, &movies); err != nil {
		return fmt.Errorf("failed to unmarshal actor seed data: %w", err)
	}
	for _, movie := range movies {
		res, err := db.Exec(query,
			movie.Title,
			movie.ReleaseYear,
			movie.Duration,
			movie.Overview,
			movie.OriginalLanguage,
		)
		if err != nil {
			return fmt.Errorf("failed to insert movie %s: %w", movie.Title, err)
		}
		movieID, resErr := res.LastInsertId()
		if resErr != nil {
			return fmt.Errorf("failed to retreive last inserted movie id: %w", resErr)
		}
		for _, genreID := range movie.GenreId {
			if _, mgErr := db.Exec(mgQuery, movieID, genreID); mgErr != nil {
				return fmt.Errorf("failed to link genre %d to movie %q: %w", genreID, movie.Title, mgErr)
			}
		}
		for _, actorID := range movie.ActorId {
			if _, maErr := db.Exec(maQuery, movieID, actorID); maErr != nil {
				return fmt.Errorf("failed to link actor %d to movie %q: %w", actorID, movie.Title, maErr)
			}
		}

	}

	return nil
}
