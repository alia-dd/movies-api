package database

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"os"
	"strings"
)

var scheme = `
	CREATE TABLE IF NOT EXISTS MOVIES(
	id 					INTEGER PRIMARY KEY,
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
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	birthdate TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name, birthdate)
	);

	CREATE TABLE IF NOT EXISTS GENRES(
	id INTEGER PRIMARY KEY,
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

	CREATE TABLE IF NOT EXISTS APIKEY(
	id   INTEGER PRIMARY KEY,
	user TEXT NOT NULL UNIQUE,
	key  TEXT NOT NULL UNIQUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`

var ApiPreFix = `test_api`

type MovieSeed struct {
	Results []models.Movies `json:"results"`
}

type actorSeed struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Character  string `json:"character"`
	PersonInfo struct {
		BirthDate string `json:"birthday"`
	} `json:"Person_Info"`
}

func InitializeMovieTable(db *sql.DB) error {
	if _, err := db.Exec(scheme); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func SeedData(db *sql.DB) error {
	cx := context.Background()

	var actors []actorSeed
	var genres []*models.Genre
	var movies []models.Movies
	mr := repository.NewMovieRepository(db)
	ar := repository.NewActorRepository(db)
	gr := repository.NewGenreRepository(db)

	// this func clears the movies/genres/actors tables for seeding process
	if clearErr := clearTables(cx, db); clearErr != nil {
		return clearErr
	}

	// get json data from the actor json file and unmarshell it into
	// actor struct and then post it the the ACTORS table
	actorBody, fileErr := os.ReadFile("internal/database/data/tmdb_actors.json")
	if fileErr != nil {
		return fileErr
	}
	if err := json.Unmarshal(actorBody, &actors); err != nil {
		return fmt.Errorf("failed to unmarshal actor seed data: %w", err)
	}
<<<<<<< HEAD
	for _, actorseed := range actors {
		actor := &models.Actor{
			Id:        actorseed.ID,
			Name:      actorseed.Name,
			BirthDate: actorseed.PersonInfo.BirthDate,
		}
		if err := ar.CreateActor(actor); err != nil {
			// if strings.Contains(err.Error(), errors.ErrDuplicateKey.Error()) {
			// 	continue
			// }
			return fmt.Errorf("failed to seed actor %q: %w", actor.Name, err)
=======
	for _, actor := range actors {
		if err := ar.CreateActor(cx, actor); err != nil {
			return fmt.Errorf("failed to seed movie %q: %w", actor.Name, err)
>>>>>>> main
		}
	}

	// get json data from the genre json file and unmarshell it into
	// genre struct and then post it the the GENRES table
	genreBody, fileErr := os.ReadFile("internal/database/data/tmdb_genres.json")
	if fileErr != nil {
		return fileErr
	}
	if err := json.Unmarshal(genreBody, &genres); err != nil {
		return fmt.Errorf("failed to unmarshal genre seed data: %w", err)
	}
	for _, genre := range genres {
		if err := gr.CreateGenre(cx, genre); err != nil {
			if strings.Contains(err.Error(), errors.ErrDuplicateKey.Error()) {
				continue
			}
			return fmt.Errorf("failed to seed genre %q: %w", genre.Name, err)
		}
	}

	// get json data from the movie json file and unmarshell it into
	// movie struct and then post it the the movie table
	// and create the actors/genre connection if there are any
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

func clearTables(cx context.Context, db *sql.DB) error {
	tx, txErr := db.BeginTx(cx, nil)
	if txErr != nil {
		return errors.ErrTransactionStart
	}
	defer tx.Rollback()
	tables := []string{"movie_actor", "movie_genre", "MOVIES", "ACTORS", "GENRES"}
	for _, table := range tables {
		if _, tableErr := tx.ExecContext(cx, `DELETE FROM `+table); tableErr != nil {
			return tableErr
		}
	}
	return tx.Commit()
}

func GenerateApiKey(db *sql.DB, user *string) (string, error) {
	query := `INSERT INTO APIKEY (user, key)VALUES (?, ?)`
	//generates 32 bit random string
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	//hashes the random bit string for security
	hashKey := sha256.Sum256(key)
	//converts the hash bit to readable string
	hashHex := hex.EncodeToString(hashKey[:])

	_, err := db.Exec(query, user, hashHex)
	if err != nil {
		return "", err
	}

	// return the unhashed hep with the prefix specified for the current user
	fullKey := ApiPreFix + "_" + hex.EncodeToString(key)
	return fullKey, nil
}

func GetApiKey(db *sql.DB, key string) bool {
	var count int
	secretkey := strings.TrimPrefix(key, ApiPreFix+"_")
	secret, err := hex.DecodeString(secretkey)
	if err != nil {
		return false
	}

	query := `SELECT count(*) FROM APIKEY WHERE key = ?`

	//hashes the random bit string for security
	hashKey := sha256.Sum256(secret)
	//converts the hash bit to readable string
	hashHex := hex.EncodeToString(hashKey[:])

	apiErr := db.QueryRow(query, hashHex).Scan(&count)
	if apiErr != nil || count <= 0 {
		return false
	}
	return true
}
