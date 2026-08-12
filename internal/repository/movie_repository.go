package repository

import (
	"database/sql"
	"fmt"
	"movies-api/internal/models"
	"strconv"
	"strings"
)

type MovieRepository struct {
	DB *sql.DB
}

// meshan waxaa galaayo sql strings data
func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}
func (r *MovieRepository) Get(f models.Filter) ([]models.Movies, error) {
	query := `SELECT m.id, m.title, m.releaseYear, m.duration, m.overview, m.originalLanguage, m.GenreId, m.ActorId, m.created_at, m.updated_at FROM MOVIES m`
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

	size, page := -1, 1
	if f.Size != "" {
		if s, err := strconv.Atoi(f.Size); err == nil {
			size = s
		}
	}
	if f.Page != "" {
		page, _ = strconv.Atoi(f.Page)
	}
	if size > -1 {
		query += ` Limit ?  OFFSET ? `
		arg = append(arg, size, (page-1)*size)

	}

	rows, rowErr := r.DB.Query(query, arg...)
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	movies := []models.Movies{}
	movieRepo := NewActorRepository(r.DB)
	genreeRepo := NewGenreRepository(r.DB)
	for rows.Next() {
		movie := models.Movies{}
		err := rows.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.GenreId, &movie.ActorId, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movieGenres := []string{}
		movieActors := []string{}
		for _, actorid := range movie.ActorId {
			payload, _ := movieRepo.FindById(actorid)
			movieActors = append(movieActors, payload.Name)
		}
		for _, genreid := range movie.GenreId {
			payload, _ := genreeRepo.GetGenreByID(genreid)
			movieGenres = append(movieGenres, payload.Name)
		}
		/// moviega waa int id hadi arabo id lee udiipi karaa string masiing kari :(
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetById(id int) (*models.Movies, error) {
	row := r.DB.QueryRow("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE id=? ", id)
	movie := models.Movies{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}
func (r *MovieRepository) GetByTitle(title string) (*models.Movies, error) {
	row := r.DB.QueryRow("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title=? ", title)
	movie := models.Movies{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) SearchByTitle(title string) ([]models.Movies, error) {
	row, rowErr := r.DB.Query("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title LIKE ?", "%"+title+"%")
	if rowErr != nil {
		return nil, rowErr
	}
	defer row.Close()
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

func (r *MovieRepository) Post(m models.Movies) (int64, error) {
	query := ` INSERT INTO MOVIES (Title, ReleaseYear, Duration, Overview, OriginalLanguage)VALUES (?, ?, ?, ?, ?)`
	mgQuery := `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
	maQuery := `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`
	res, postErr := r.DB.Exec(query, m.Title, m.ReleaseYear, m.Duration, m.Overview, m.OriginalLanguage)
	if postErr != nil {
		return 0, fmt.Errorf("failed to insert MOVIES into table: %w", postErr)
	}
	movieID, resErr := res.LastInsertId()
	if resErr != nil {
		return 0, fmt.Errorf("failed to retreive last inserted MOVIES id: %w", resErr)
	}
	for _, genreID := range m.GenreId {
		if _, mgErr := r.DB.Exec(mgQuery, movieID, genreID); mgErr != nil {
			return 0, fmt.Errorf("failed to link genre %d to MOVIES %q: %w", genreID, m.Title, mgErr)
		}
	}
	// var a ActorsRepository
	for _, actorID := range m.ActorId {

		if _, maErr := r.DB.Exec(maQuery, movieID, actorID); maErr != nil {
			return 0, fmt.Errorf("failed to link actor %d to MOVIES %q: %w", actorID, m.Title, maErr)
		}
	}
	return movieID, nil
}

func (r *MovieRepository) Patch() {
}

func (r *MovieRepository) Delete(id int, force bool) (int64, error) {
	if force {
		m, payloadErr := r.GetById(id)
		if payloadErr != nil {
			return 0, payloadErr
		}
		mgQuery := `DELETE FROM movie_genre WHERE movieId = ? AND genreId = ?`
		maQuery := `DELETE FROM movie_actor WHERE movieId = ? AND actorId = ?`
		for _, genreID := range m.GenreId {
			if _, mgErr := r.DB.Exec(mgQuery, id, genreID); mgErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %q: %w", genreID, m.Title, mgErr)
			}
		}
		for _, actorID := range m.ActorId {
			if _, maErr := r.DB.Exec(maQuery, id, actorID); maErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %q: %w", actorID, m.Title, maErr)
			}
		}
	}
	rows, deleteErr := r.DB.Exec("DELETE FROM MOVIES WHERE id=?", id)
	if deleteErr != nil {
		return 0, deleteErr
	}
	affectedR, _ := rows.RowsAffected()
	return affectedR, nil
}
