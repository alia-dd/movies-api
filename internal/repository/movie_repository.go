package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	// "fmt"
	"movies-api/internal/models"
	"strconv"
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (r *MovieRepository) Get(cx context.Context, f *models.Filter) ([]models.MoviesDisplay, error) {

	query := `SELECT m.id, m.title, m.releaseYear, m.duration, m.overview, m.originalLanguage, m.created_at, m.updated_at FROM MOVIES m`
	extraQuery := ``
	arg := []any{}

	if f.Genre != "" {
		extraQuery += ` SELECT mg.movieId FROM movie_genre mg WHERE mg.genreId IN (SELECT g.id FROM GENRES g WHERE (g.id = ? or g.name = ?)) `
		arg = append(arg, f.Genre, f.Genre)
	}
	if f.Actor != "" && f.Genre != "" {
		extraQuery += ` INTERSECT`
	}
	if f.Actor != "" {
		extraQuery += ` SELECT ma.movieId FROM movie_actor ma WHERE ma.actorId IN (SELECT a.id FROM ACTORS a WHERE (a.id = ? or a.name = ?))`
		arg = append(arg, f.Actor, f.Actor)
	}
	if f.Actor != "" || f.Genre != "" {
		extraQuery = ` WHERE m.id IN (` + extraQuery + `)`
	}
	if f.Year != "" {
		if f.Actor != "" || f.Genre != "" {
			extraQuery += ` AND `
		}
		extraQuery += ` m.releaseYear = ?`
		arg = append(arg, f.Year)
	}
	query += extraQuery
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

	rows, rowErr := r.DB.QueryContext(cx, query, arg...)
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	movies := []models.MoviesDisplay{}
	for rows.Next() {
		movie := models.Movies{}
		err := rows.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}

		duration := formatDurations(movie.Duration)
		movies = append(movies, models.MoviesDisplay{
			Id:               movie.Id,
			Title:            movie.Title,
			ReleaseYear:      movie.ReleaseYear,
			Duration:         duration,
			Overview:         movie.Overview,
			OriginalLanguage: movie.OriginalLanguage,
			CreatedAt:        movie.CreatedAt,
			UpdatedAt:        movie.UpdatedAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range movies {
		movieActors, movieGenres := r.getInfo(cx, movies[i].Id)
		movies[i].Actors = movieActors
		movies[i].Genres = movieGenres
	}

	return movies, nil
}

func (r *MovieRepository) GetById(cx context.Context, id int) (*models.MoviesDisplay, error) {
	row := r.DB.QueryRowContext(cx, "SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE id=? ", id)
	movie := models.MoviesDisplay{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	movieActors, movieGenres := r.getInfo(cx, movie.Id)
	movie.Actors = movieActors
	movie.Genres = movieGenres
	return &movie, nil
}
func (r *MovieRepository) GetByTitle(cx context.Context, title string) (*models.MoviesDisplay, error) {
	row := r.DB.QueryRowContext(cx, "SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title=? ", title)
	movie := models.MoviesDisplay{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	movieActors, movieGenres := r.getInfo(cx, movie.Id)
	movie.Actors = movieActors
	movie.Genres = movieGenres
	return &movie, nil
}

func (r *MovieRepository) SearchByTitle(cx context.Context, title string) ([]models.MoviesDisplay, error) {
	row, rowErr := r.DB.QueryContext(cx, "SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title LIKE ? ", "%"+title+"%")
	if rowErr != nil {
		return nil, rowErr
	}
	defer row.Close()
	movies := []models.MoviesDisplay{}
	for row.Next() {
		movie := models.MoviesDisplay{}
		err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	for i := range movies {
		movieActors, movieGenres := r.getInfo(cx, movies[i].Id)
		movies[i].Actors = movieActors
		movies[i].Genres = movieGenres
	}

	return movies, nil
}

func (r *MovieRepository) Post(cx context.Context, m models.Movies) (int64, error) {
	actorRepo := NewActorRepository(r.DB)
	genreRepo := NewGenreRepository(r.DB)
	tx, txErr := r.DB.BeginTx(cx, nil)
	if txErr != nil {
		return 0, fmt.Errorf("Failed to begin transaction: %w", txErr)
	}
	query := ` INSERT INTO MOVIES (Title, ReleaseYear, Duration, Overview, OriginalLanguage)VALUES (?, ?, ?, ?, ?)`
	mgQuery := `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
	maQuery := `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`
	defer tx.Rollback()
	res, postErr := tx.ExecContext(cx, query, m.Title, m.ReleaseYear, m.Duration, m.Overview, m.OriginalLanguage)
	if postErr != nil {
		return 0, fmt.Errorf("failed to insert MOVIES into table: %w", postErr)
	}
	movieID, resErr := res.LastInsertId()
	if resErr != nil {
		return 0, fmt.Errorf("failed to retreive last inserted MOVIES id: %w", resErr)
	}
	for _, genreID := range m.GenreId {
		genre, err := genreRepo.GetGenreByID(genreID)
		if err != nil || genre == nil {
			continue
		}
		if _, mgErr := tx.ExecContext(cx, mgQuery, movieID, genreID); mgErr != nil {
			return 0, fmt.Errorf("failed to link genre %d to MOVIES %q: %w", genreID, m.Title, mgErr)
		}
	}

	for _, actorID := range m.ActorId {
		actor, err := actorRepo.FindById(actorID)
		if err != nil || actor == nil {
			continue
		}
		if _, maErr := tx.ExecContext(cx, maQuery, movieID, actorID); maErr != nil {
			return 0, fmt.Errorf("failed to link actor %d to MOVIES %q: %w", actorID, m.Title, maErr)
		}
	}
	tx.Commit()
	return movieID, nil
}

func (r *MovieRepository) Patch(cx context.Context, id int, m models.MovieUpdate) error {
	mgQuery := `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
	maQuery := `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`

	mgDeleteQuery := `DELETE FROM movie_genre WHERE movieId = ? AND genreId = ?`
	maDeleteQuery := `DELETE FROM movie_actor WHERE movieId = ? AND actorId = ?`

	tx, txErr := r.DB.BeginTx(cx, nil)
	if txErr != nil {
		return fmt.Errorf("Failed to begin transaction: %w", txErr)
	}
	query := `UPDATE MOVIES SET updated_at = ? `
	now := time.Now()
	extraQuery := []string{}
	arg := []any{now}
	if m.Title != nil {
		extraQuery = append(extraQuery, ` title = ? `)
		arg = append(arg, *m.Title)
	}
	if m.ReleaseYear != nil {
		extraQuery = append(extraQuery, ` releaseYear = ? `)
		arg = append(arg, *m.ReleaseYear)
	}
	if m.Duration != nil {
		extraQuery = append(extraQuery, ` duration = ?`)
		arg = append(arg, *m.Duration)
	}
	if m.Overview != nil {
		extraQuery = append(extraQuery, ` Overview = ?`)
		arg = append(arg, *m.Overview)
	}
	if m.OriginalLanguage != nil {
		extraQuery = append(extraQuery, ` OriginalLanguage = ? `)
		arg = append(arg, *m.OriginalLanguage)
	}
	defer tx.Rollback()
	if len(extraQuery) != 0 {
		query = fmt.Sprintf("%s, %s WHERE id = ? ", query, strings.Join(extraQuery, ", "))
		arg = append(arg, id)
		_, patchErr := tx.ExecContext(cx, query, arg...)
		if patchErr != nil {
			return fmt.Errorf("failed to Update MOVIE id: %d %w", id, patchErr)
		}
	}

	if len(m.AddActorIDs) > 0 {
		for _, actorID := range m.AddActorIDs {
			if _, maErr := tx.ExecContext(cx, maQuery, id, actorID); maErr != nil {
				if strings.Contains(maErr.Error(), "UNIQUE constraint failed") {
					continue
				}
				return fmt.Errorf("failed to link actor %d to MOVIES %d: %w", actorID, id, maErr)
			}
		}
	}
	if len(m.RemoveActorIDs) > 0 {
		for _, actorID := range m.RemoveActorIDs {
			if _, maErr := tx.ExecContext(cx, maDeleteQuery, id, actorID); maErr != nil {
				return fmt.Errorf("failed to remove link actor %d to MOVIES %d: %w", actorID, id, maErr)
			}
		}
	}
	if len(m.AddGenreIDs) > 0 {
		for _, genreID := range m.AddGenreIDs {
			if _, mgErr := tx.ExecContext(cx, mgQuery, id, genreID); mgErr != nil {
				if strings.Contains(mgErr.Error(), "UNIQUE constraint failed") {
					continue
				}
				return fmt.Errorf("failed to link genre %d to MOVIES %d: %w", genreID, id, mgErr)
			}
		}
	}
	if len(m.RemoveGenreIDs) > 0 {
		for _, genreID := range m.RemoveGenreIDs {
			if _, mgErr := tx.ExecContext(cx, mgDeleteQuery, id, genreID); mgErr != nil {
				return fmt.Errorf("failed to remove link genre %d to MOVIES %d: %w", genreID, id, mgErr)
			}
		}
	}
	return tx.Commit()
}

func (r *MovieRepository) Delete(cx context.Context, MovieId int, force bool) (int64, error) {
	if force {
		m, payloadErr := r.GetById(cx, MovieId)
		if payloadErr != nil {
			return 0, payloadErr
		}
		mgQuery := `DELETE FROM movie_genre WHERE movieId = ? AND genreId = ?`
		maQuery := `DELETE FROM movie_actor WHERE movieId = ? AND actorId = ?`
		genresId, _ := r.GetMovie_genre(cx, MovieId)
		for _, genreID := range genresId {
			if _, mgErr := r.DB.ExecContext(cx, mgQuery, MovieId, genreID); mgErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %s: %w", genreID, m.Title, mgErr)
			}
		}
		actorsId, _ := r.GetMovie_actor(cx, MovieId)
		for _, actorID := range actorsId {
			if _, maErr := r.DB.ExecContext(cx, maQuery, MovieId, actorID); maErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %s: %w", actorID, m.Title, maErr)
			}
		}
	}
	rows, deleteErr := r.DB.ExecContext(cx, "DELETE FROM MOVIES WHERE id=?", MovieId)
	if deleteErr != nil {
		return 0, deleteErr
	}
	affectedR, _ := rows.RowsAffected()
	return affectedR, nil
}

func (r *MovieRepository) GetMovie_genre(cx context.Context, MovieId int) ([]int, int) {
	count := 0
	genresId := []int{}
	mgQuery := `SELECT genreId FROM movie_genre WHERE movieId = ?`
	rows, err := r.DB.QueryContext(cx, mgQuery, MovieId)
	if err != nil {
		return []int{}, 0
	}
	defer rows.Close()
	for rows.Next() {
		var genreId int
		rows.Scan(&genreId)
		genresId = append(genresId, genreId)
		count++
	}
	if err := rows.Err(); err != nil {
		return []int{}, 0
	}
	return genresId, count
}

func (r *MovieRepository) GetMovie_actor(cx context.Context, MovieId int) ([]int, int) {
	count := 0
	actorsId := []int{}
	maQuery := `SELECT actorId FROM movie_actor WHERE movieId = ?`

	rows, err := r.DB.QueryContext(cx, maQuery, MovieId)
	if err != nil {
		return []int{}, 0
	}
	defer rows.Close()
	for rows.Next() {
		var actorId int
		rows.Scan(&actorId)
		actorsId = append(actorsId, actorId)
		count++
	}
	if err := rows.Err(); err != nil {
		return []int{}, 0
	}
	return actorsId, count
}

func (r *MovieRepository) getInfo(cx context.Context, MovieId int) ([]string, []string) {
	movieGenres := []string{}
	movieActors := []string{}

	actorRepo := NewActorRepository(r.DB)
	genreRepo := NewGenreRepository(r.DB)

	actorsId, _ := r.GetMovie_actor(cx, MovieId)
	genresId, _ := r.GetMovie_genre(cx, MovieId)

	for _, actorid := range actorsId {
		actor, err := actorRepo.FindById(actorid)
		if err != nil {
			continue
		}
		movieActors = append(movieActors, actor.Name)
	}
	for _, genreId := range genresId {
		genre, err := genreRepo.GetGenreByID(genreId)
		if err != nil {
			continue
		}
		movieGenres = append(movieGenres, genre.Name)
	}
	return movieActors, movieGenres
}

func formatDurations(duration uint16) string {
	return ""
}
