package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"os"
)

var scheme = `
	CREATE TABLE IF NOT EXISTS MOVIES(
	id 					INTEGER PRIMARY KEY AUTOINCREMENT,
	title       		TEXT NOT NULL,
	releaseYear 		INTEGER NOT NULL,
	duration    		INTEGER,
	Overview 			TEXT,
	OriginalLanguage 	TEXT,
	created_at 			DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at 			DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(title,releaseYear,duration)
	);

	CREATE TABLE IF NOT EXISTS ACTORS(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	birthdate TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS GENRES(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS movie_genre(
	genreId      		INTEGER NOT NULL,
	movieId      		INTEGER NOT NULL,
	PRIMARY KEY (movieId, genreId),
    FOREIGN KEY (movieId) REFERENCES MOVIES (id) ,
	FOREIGN KEY (genreId) REFERENCES GENRES (id) 
	);

	CREATE TABLE IF NOT EXISTS movie_actor(
	actorId      		INTEGER NOT NULL,
	movieId      		INTEGER NOT NULL,
	PRIMARY KEY (movieId, actorId),
    FOREIGN KEY (movieId) REFERENCES MOVIES (id) ,
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

	if _, err := db.Exec(scheme); err != nil {
		return fmt.Errorf("failed to create MOVIES table: %w", err)
	}

	return nil
}
func tableExist(db *sql.DB) bool {
	var count int
	query := `SELECT count(*) FROM sqlite_master WHERE name ='MOVIES' and type='table'`
	err := db.QueryRow(query).Scan(&count)
	if err != nil || count <= 0 {
		return false
	}
	return true
}

func SeedData(db *sql.DB) error {
	cx := context.Background()

	var actors []*models.Actor
	var genres []*models.Genre
	var movies []models.Movies
	mr := repository.NewMovieRepository(db)
	ar := repository.NewActorRepository(db)
	gr := repository.NewGenreRepository(db)

	actorBody, fileErr := os.ReadFile("internal/database/data/tmdb_actors.json")
	if fileErr != nil {
		return fileErr
	}
	if err := json.Unmarshal(actorBody, &actors); err != nil {
		return fmt.Errorf("failed to unmarshal actor seed data: %w", err)
	}

	for _, actor := range actors {
		if err := ar.CreateActor(actor); err != nil {
			return fmt.Errorf("failed to seed movie %q: %w", actor.Name, err)
		}
	}

	genreBody, fileErr := os.ReadFile("internal/database/data/tmdb_genres.json")
	if fileErr != nil {
		return fileErr
	}
	if err := json.Unmarshal(genreBody, &genres); err != nil {
		return fmt.Errorf("failed to unmarshal genre seed data: %w", err)
	}

	for _, genre := range genres {
		if err := gr.CreateGenre(genre); err != nil {
			return fmt.Errorf("failed to seed genre %q: %w", genre.Name, err)
		}
	}

	movieBody, fileErr := os.ReadFile("internal/database/data/tmdb_movies.json")
	if fileErr != nil {
		return fileErr
	}
	if err := json.Unmarshal(movieBody, &movies); err != nil {
		return fmt.Errorf("failed to unmarshal movie seed data: %w", err)
	}
	for _, movie := range movies {
		if _, err := mr.Post(cx, movie); err != nil {
			return fmt.Errorf("failed to seed movie %q: %w", movie.Title, err)
		}
	}
	return nil
}
