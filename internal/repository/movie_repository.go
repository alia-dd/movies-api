package repository

import (
	"database/sql"
	"fmt"
	"movies-api/internal/models"
	"strconv"
	"strings"
)

type DatabaseConnection struct {
	DB *sql.DB
}

// meshan waxaa galaayo sql strings data
func NewdbConnection(db *sql.DB) *DatabaseConnection {
	return &DatabaseConnection{DB: db}
}
func (c *DatabaseConnection) Get(f models.Filter) ([]models.Movies, error) {
	query := `SELECT m.id, m.title, m.releaseYear, m.duration, m.overview, m.originalLanguage, m.created_at, m.updated_at FROM movie m`
	extraQuery := []string{}
	arg := []any{}

	if f.Genre != "" {
		query += ` JOIN movie_genre mg ON m.id = mg.movieId `
		extraQuery = append(extraQuery, ` mg.genreId=?`)
		arg = append(arg, f.Genre)
	}
	if f.Actor != "" {
		query += ` JOIN movie_actor ma ON m.id  = ma.movieId`
		extraQuery = append(extraQuery, `ma.actorId = ? `)
		arg = append(arg, f.Actor)
	}
	if f.Year != "" {
		extraQuery = append(extraQuery, ` m.releaseYear = ? `)
		arg = append(arg, f.Year)
	}
	if len(extraQuery) > 0 {
		query += " WHERE " + strings.Join(extraQuery, " AND ")
	}

	size, page := 0, 0
	if f.Size != "" {
		size, _ = strconv.Atoi(f.Size)
		query += `Limit = ? `
		arg = append(arg, size)
	}
	if f.Page != "" {
		page, _ = strconv.Atoi(f.Page)
		query += `OFFSET = ? `
		arg = append(arg, page*size)
	}

	rows, rowErr := c.DB.Query(query, arg...)
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	movies := []models.Movies{}

	for rows.Next() {
		movie := models.Movies{}
		err := rows.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (c *DatabaseConnection) GetById(id int) (*models.Movies, error) {
	row := c.DB.QueryRow("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM movie WHERE id=? ", id)
	movie := models.Movies{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}
func (c *DatabaseConnection) GetByTitle(actorName string) (*models.Movies, error) {
	row := c.DB.QueryRow("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM movie WHERE title=? ", actorName)
	movie := models.Movies{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (c *DatabaseConnection) SearchByTitle(actorName string) ([]models.Movies, error) {
	row, rowErr := c.DB.Query("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM movie WHERE title LIKE ?", "%"+actorName+"%")
	if rowErr != nil {
		return nil, rowErr
	}
	movies := []models.Movies{}
	for row.Next() {
		movie := models.Movies{}
		err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

func (c *DatabaseConnection) Post(m models.Movies) error {
	query := ` INSERT INTO movie (Title, ReleaseYear, Duration, Overview, OriginalLanguage)VALUES (?, ?, ?, ?, ?)`
	mgQuery := `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
	maQuery := `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`
	res, postErr := c.DB.Exec(query, m.Title, m.ReleaseYear, m.Duration, m.Overview, m.OriginalLanguage)
	if postErr != nil {
		return fmt.Errorf("failed to insert movie into table: %w", postErr)
	}
	movieID, resErr := res.LastInsertId()
	if resErr != nil {
		return fmt.Errorf("failed to retreive last inserted movie id: %w", resErr)
	}
	for _, genreID := range m.GenreId {
		if _, mgErr := c.DB.Exec(mgQuery, movieID, genreID); mgErr != nil {
			return fmt.Errorf("failed to link genre %d to movie %q: %w", genreID, m.Title, mgErr)
		}
	}
	for _, actorID := range m.ActorId {
		if _, maErr := c.DB.Exec(maQuery, movieID, actorID); maErr != nil {
			return fmt.Errorf("failed to link actor %d to movie %q: %w", actorID, m.Title, maErr)
		}
	}
	return nil
}

func (c *DatabaseConnection) Patch() {
}

func (c *DatabaseConnection) Delete(id int, force bool) (int64, error) {
	if force {
		m, payloadErr := c.GetById(id)
		if payloadErr != nil {
			return 0, payloadErr
		}
		// err := fmt.Sprintf("Cannot delete Movie %s because it has %d associated actors and %d associated genres.", payload.Title, len(payload.ActorId), len(payload.GenreId))
		// return 0, errors.New(err)

		mgQuery := `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
		maQuery := `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`
		for _, genreID := range m.GenreId {
			if _, mgErr := c.DB.Exec(mgQuery, id, genreID); mgErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and movie %q: %w", genreID, m.Title, mgErr)
			}
		}
		for _, actorID := range m.ActorId {
			if _, maErr := c.DB.Exec(maQuery, id, actorID); maErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and movie %q: %w", actorID, m.Title, maErr)
			}
		}
	}
	rows, deleteErr := c.DB.Exec("DELETE FROM movie WHERE id=?", id)
	if deleteErr != nil {
		return 0, deleteErr
	}
	affectedR, _ := rows.RowsAffected()
	return affectedR, nil
}
