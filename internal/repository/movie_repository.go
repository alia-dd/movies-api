package repository

import (
	"context"
	"database/sql"
	"fmt"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"strconv"
	"strings"
	"time"
)

const (
	getAllSelectQuery          = `SELECT m.id, m.title, m.releaseYear, m.duration, m.overview, m.originalLanguage, m.created_at, m.updated_at FROM MOVIES m`
	getAllSelectGenreQuery     = ` SELECT mg.movieId FROM movie_genre mg WHERE mg.genreId IN (SELECT g.id FROM GENRES g WHERE (g.id = ? or g.name = ?)) `
	getAllSelectActorQuery     = ` SELECT ma.movieId FROM movie_actor ma WHERE ma.actorId IN (SELECT a.id FROM ACTORS a WHERE (a.id = ? or a.name = ?))`
	getAllWithLimitOffsetQuery = ` Limit ?  OFFSET ? `
	getByIdQuery               = `SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE id=? `
	getByTitleQuery            = `SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title = ? `
	SearchByTitleQuery         = `SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title LIKE ? `
	getAllSelectCountQuery     = `SELECT COUNT(*) FROM MOVIES `

	insertMovieQuery      = `INSERT INTO MOVIES (Title, ReleaseYear, Duration, Overview, OriginalLanguage)VALUES (?, ?, ?, ?, ?)`
	inserMovieGenreQuery  = `INSERT INTO movie_genre (movieId, genreId) VALUES (?, ?)`
	insertMovieActorQuery = `INSERT INTO movie_actor (movieId, actorId) VALUES (?, ?)`

	deleteMovieQuery                = `DELETE FROM MOVIES WHERE id=?`
	deleteMovieGenreConnectionQuery = `DELETE FROM movie_genre WHERE movieId = ? AND genreId = ?`
	deleteMovieActorConnectionQuery = `DELETE FROM movie_actor WHERE movieId = ? AND actorId = ?`

	getGenresIdBymovieIdQuery = `SELECT genreId FROM movie_genre WHERE movieId = ?`
	getActorsIdBymovieIdQuery = `SELECT actorId FROM movie_actor WHERE movieId = ?`

	patchMovieById = `UPDATE MOVIES SET updated_at = ? `
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (r *MovieRepository) Get(cx context.Context, f *models.Filter) ([]models.MoviesDisplay, int, error) {
	total := 0
	err := r.DB.QueryRowContext(cx, getAllSelectCountQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := getAllSelectQuery
	extraQuery := ``
	arg := []any{}

	if f.Genre != "" {
		extraQuery += getAllSelectGenreQuery
		arg = append(arg, f.Genre, f.Genre)
	}
	if f.Actor != "" && f.Genre != "" {
		extraQuery += ` INTERSECT`
	}
	if f.Actor != "" {
		extraQuery += getAllSelectActorQuery
		arg = append(arg, f.Actor, f.Actor)
	}
	if f.Actor != "" || f.Genre != "" {
		extraQuery = ` m.id IN (` + extraQuery + `)`
	}
	if f.Year != "" {
		if f.Actor != "" || f.Genre != "" {
			extraQuery += ` AND `
		}
		extraQuery += ` m.releaseYear = ?`
		arg = append(arg, f.Year)
	}
	if extraQuery != "" {
		extraQuery = ` WHERE ` + extraQuery
		query += extraQuery
	}
	size, page := -1, 0
	if f.Size != "" {
		if s, err := strconv.Atoi(f.Size); err == nil {
			size = s
			if size > total {
				return nil, 0, errors.ErrInvalidInput
			}
		}
	}

	if f.Page != "" {
		page, _ = strconv.Atoi(f.Page)
		f.Page = strconv.Itoa(page)
	}
	if size > -1 {
		query += getAllWithLimitOffsetQuery
		arg = append(arg, size, page*size)

	} else {
		f.Size = strconv.Itoa(total)
	}

	rows, rowErr := r.DB.QueryContext(cx, query, arg...)
	if rowErr != nil {
		return nil, 0, rowErr
	}
	defer rows.Close()

	movies := []models.MoviesDisplay{}
	for rows.Next() {
		movie := models.Movies{}
		err := rows.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, 0, err
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
		return nil, 0, err
	}

	for i := range movies {
		movieActors, _ := r.GetActorsForMovie(cx, movies[i].Id)
		movieGenres, _ := r.GetGenresForMovie(cx, movies[i].Id)
		movies[i].Actors = movieActors
		movies[i].Genres = movieGenres
	}

	return movies, total, nil
}

func (r *MovieRepository) GetById(cx context.Context, id int) (*models.MoviesDisplay, error) {
	row := r.DB.QueryRowContext(cx, getByIdQuery, id)
	movie := models.MoviesDisplay{}
	var duration uint16
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	movieActors, _ := r.GetActorsForMovie(cx, movie.Id)
	movieGenres, _ := r.GetGenresForMovie(cx, movie.Id)
	movie.Actors = movieActors
	movie.Genres = movieGenres
	movie.Duration = formatDurations(duration)
	return &movie, nil
}

func (r *MovieRepository) GetByTitle(cx context.Context, title string) (*models.MoviesDisplay, error) {
	row := r.DB.QueryRowContext(cx, getByTitleQuery, title)
	movie := models.MoviesDisplay{}
	var duration uint16
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	movieActors, _ := r.GetActorsForMovie(cx, movie.Id)
	movieGenres, _ := r.GetGenresForMovie(cx, movie.Id)
	movie.Actors = movieActors
	movie.Genres = movieGenres
	movie.Duration = formatDurations(duration)
	return &movie, nil
}

func (r *MovieRepository) SearchByTitle(cx context.Context, title string) ([]models.MoviesDisplay, error) {

	row, rowErr := r.DB.QueryContext(cx, SearchByTitleQuery, "%"+title+"%")
	if rowErr != nil {
		return nil, rowErr
	}
	defer row.Close()
	movies := []models.MoviesDisplay{}
	for row.Next() {
		movie := models.MoviesDisplay{}
		var duration uint16
		err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movie.Duration = formatDurations(duration)
		movies = append(movies, movie)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	for i := range movies {
		movieActors, _ := r.GetActorsForMovie(cx, movies[i].Id)
		movieGenres, _ := r.GetGenresForMovie(cx, movies[i].Id)
		movies[i].Actors = movieActors
		movies[i].Genres = movieGenres
	}

	return movies, nil
}

func (r *MovieRepository) Post(cx context.Context, m models.Movies) (int64, error) {
	tx, txErr := r.DB.BeginTx(cx, nil)
	if txErr != nil {
		return 0, errors.ErrTransactionStart
	}

	defer tx.Rollback() // if there is error revert changes to the db back to before the changes
	res, postErr := tx.ExecContext(cx, insertMovieQuery, m.Title, m.ReleaseYear, m.Duration, m.Overview, m.OriginalLanguage)
	if postErr != nil {
		return 0, fmt.Errorf("failed to insert movie into table: %w", postErr)
	}
	movieID, resErr := res.LastInsertId()
	if resErr != nil {
		return 0, fmt.Errorf("failed to retreive last inserted movie id: %w", resErr)
	}
	for _, genreID := range m.GenreId {
		if _, mgErr := tx.ExecContext(cx, inserMovieGenreQuery, movieID, genreID); mgErr != nil {
			if strings.Contains(mgErr.Error(), "UNIQUE constraint failed") {
				continue
			}
			return 0, fmt.Errorf("failed to link genre %d to MOVIES %q: %w", genreID, m.Title, mgErr)
		}
	}

	for _, actorID := range m.ActorId {
		if _, maErr := tx.ExecContext(cx, insertMovieActorQuery, movieID, actorID); maErr != nil {
			if strings.Contains(maErr.Error(), "UNIQUE constraint failed") {
				continue
			}
			return 0, fmt.Errorf("failed to link actor %d to MOVIES %q: %w", actorID, m.Title, maErr)
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, errors.ErrTransactionCommit
	}
	return movieID, nil
}

func (r *MovieRepository) Patch(cx context.Context, id int, m models.MovieUpdate) error {
	tx, txErr := r.DB.BeginTx(cx, nil)
	if txErr != nil {
		return errors.ErrTransactionStart
	}

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
	defer tx.Rollback() // if there is error revert changes to the db back to before the changes
	if len(extraQuery) != 0 {
		query := fmt.Sprintf("%s, %s WHERE id = ? ", patchMovieById, strings.Join(extraQuery, ", "))
		arg = append(arg, id)
		_, patchErr := tx.ExecContext(cx, query, arg...)
		if patchErr != nil {
			return fmt.Errorf("failed to Update MOVIE id: %d %w", id, patchErr)
		}
	} else {
		_, patchErr := tx.ExecContext(cx, patchMovieById, arg...)
		if patchErr != nil {
			return fmt.Errorf("failed to Update MOVIE id: %d %w", id, patchErr)
		}
	}

	if len(m.AddActorIDs) > 0 {
		for _, actorID := range m.AddActorIDs {
			if _, maErr := tx.ExecContext(cx, insertMovieActorQuery, id, actorID); maErr != nil {
				if strings.Contains(maErr.Error(), "UNIQUE constraint failed") {
					continue
				}
				return fmt.Errorf("failed to link actor %d to MOVIES %d: %w", actorID, id, maErr)
			}
		}
	}
	if len(m.RemoveActorIDs) > 0 {
		for _, actorID := range m.RemoveActorIDs {
			if _, maErr := tx.ExecContext(cx, deleteMovieActorConnectionQuery, id, actorID); maErr != nil {
				return fmt.Errorf("failed to remove link actor %d to MOVIES %d: %w", actorID, id, maErr)
			}
		}
	}
	if len(m.AddGenreIDs) > 0 {
		for _, genreID := range m.AddGenreIDs {
			if _, mgErr := tx.ExecContext(cx, inserMovieGenreQuery, id, genreID); mgErr != nil {
				if strings.Contains(mgErr.Error(), "UNIQUE constraint failed") {
					continue
				}
				return fmt.Errorf("failed to link genre %d to MOVIES %d: %w", genreID, id, mgErr)
			}
		}
	}
	if len(m.RemoveGenreIDs) > 0 {
		for _, genreID := range m.RemoveGenreIDs {
			if _, mgErr := tx.ExecContext(cx, deleteMovieGenreConnectionQuery, id, genreID); mgErr != nil {
				return fmt.Errorf("failed to remove link genre %d to MOVIES %d: %w", genreID, id, mgErr)
			}
		}
	}
	return tx.Commit()
}

func (r *MovieRepository) Delete(cx context.Context, MovieId int, force bool) (int64, error) {
	var genresId, actorsId []int
	var movieTitle string

	if force {
		m, payloadErr := r.GetById(cx, MovieId)
		if payloadErr != nil {

			return 0, payloadErr
		}
		movieTitle = m.Title
		genresId, _ = r.GetMovie_genre(cx, MovieId)
		actorsId, _ = r.GetMovie_actor(cx, MovieId)
	}
	tx, txErr := r.DB.BeginTx(cx, nil)
	if txErr != nil {
		return 0, errors.ErrTransactionStart
	}
	defer tx.Rollback()

	if force {
		for _, genreID := range genresId {
			if _, mgErr := tx.ExecContext(cx, deleteMovieGenreConnectionQuery, MovieId, genreID); mgErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %s: %w", genreID, movieTitle, mgErr)
			}
		}

		for _, actorID := range actorsId {
			if _, maErr := tx.ExecContext(cx, deleteMovieActorConnectionQuery, MovieId, actorID); maErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw actor %d and MOVIES %s: %w", actorID, movieTitle, maErr)
			}
		}
	}
	rows, deleteErr := tx.ExecContext(cx, deleteMovieQuery, MovieId)
	if deleteErr != nil {
		return 0, deleteErr
	}
	affectedR, _ := rows.RowsAffected()

	if err := tx.Commit(); err != nil {
		return 0, errors.ErrTransactionCommit
	}
	return affectedR, nil
}

// this method returns the genreIds of provided movieId
func (r *MovieRepository) GetMovie_genre(cx context.Context, MovieId int) ([]int, int) {
	count := 0
	genresId := []int{}
	rows, err := r.DB.QueryContext(cx, getGenresIdBymovieIdQuery, MovieId)
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

// this method returns the actorIds of provided movieId
func (r *MovieRepository) GetMovie_actor(cx context.Context, MovieId int) ([]int, int) {
	count := 0
	actorsId := []int{}
	rows, err := r.DB.QueryContext(cx, getActorsIdBymovieIdQuery, MovieId)
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

// this method returnes the names of genres of a provided movieId by using the GetMovie_genre function
func (r *MovieRepository) GetGenresForMovie(cx context.Context, MovieId int) ([]string, error) {

	movieGenres := []string{}
	genreRepo := NewGenreRepository(r.DB)
	genresId, _ := r.GetMovie_genre(cx, MovieId)

	for _, genreId := range genresId {
		genre, err := genreRepo.GetGenreByID(cx, genreId)
		if err != nil {
			continue
		}
		movieGenres = append(movieGenres, genre.Name)
	}
	return movieGenres, nil
}

// this method returnes the names of actors of a provided movieId by using the GetMovie_actor function
func (r *MovieRepository) GetActorsForMovie(cx context.Context, MovieId int) ([]string, error) {

	movieActors := []string{}
	actorRepo := NewActorRepository(r.DB)
	actorsId, _ := r.GetMovie_actor(cx, MovieId)

	for _, actorid := range actorsId {
		actor, err := actorRepo.FindById(cx, actorid)
		if err != nil {
			continue
		}
		movieActors = append(movieActors, actor.Name)
	}
	return movieActors, nil
}

func formatDurations(duration uint16) string {
	return fmt.Sprintf("%02d h %02d m", duration/60, duration%60)
}
