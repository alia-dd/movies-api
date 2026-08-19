package repository

import (
	"database/sql"
	"fmt"

	// "fmt"
	"movies-api/internal/models"
	"strconv"
)

type MovieRepository struct {
	DB *sql.DB
}

// meshan waxaa galaayo sql strings data
func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (r *MovieRepository) Get(f *models.Filter) ([]models.MoviesDisplay, error) {

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

	rows, rowErr := r.DB.Query(query, arg...)
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
		movieActors, _ := r.GetActorsForMovie(movies[i].Id)
		movieGenres, _ := r.GetGenresForMovie(movies[i].Id)
		movies[i].Actors = movieActors
		movies[i].Genres = movieGenres
	}

	return movies, nil
}

func (r *MovieRepository) GetById(id int) (*models.MoviesDisplay, error) {
	row := r.DB.QueryRow("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE id=? ", id)
	movie := models.MoviesDisplay{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	movieActors, _ := r.GetActorsForMovie(movie.Id)
	movieGenres, _ := r.GetGenresForMovie(movie.Id)
	movie.Actors = movieActors
	movie.Genres = movieGenres
	return &movie, nil
}
func (r *MovieRepository) GetByTitle(title string) (*models.MoviesDisplay, error) {
	row := r.DB.QueryRow("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title=? ", title)
	movie := models.MoviesDisplay{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.Overview, &movie.OriginalLanguage, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	movieActors, _ := r.GetActorsForMovie(movie.Id)
	movieGenres, _ := r.GetGenresForMovie(movie.Id)
	movie.Actors = movieActors
	movie.Genres = movieGenres
	return &movie, nil
}

func (r *MovieRepository) SearchByTitle(title string) ([]models.MoviesDisplay, error) {
	row, rowErr := r.DB.Query("SELECT id, title, releaseYear, duration, overview, originalLanguage, created_at, updated_at FROM MOVIES WHERE title LIKE ? ", "%"+title+"%")
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
		movieActors, _ := r.GetActorsForMovie(movies[i].Id)
		movieGenres, _ := r.GetGenresForMovie(movies[i].Id)
		movies[i].Actors = movieActors
		movies[i].Genres = movieGenres
	}

	return movies, nil
}

func (r *MovieRepository) Post(m models.Movies) (int64, error) {
	actorRepo := NewActorRepository(r.DB)
	genreRepo := NewGenreRepository(r.DB)

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
		genre, err := genreRepo.GetGenreByID(genreID)
		if err != nil || genre == nil {
			continue
		}
		if _, mgErr := r.DB.Exec(mgQuery, movieID, genreID); mgErr != nil {
			return 0, fmt.Errorf("failed to link genre %d to MOVIES %q: %w", genreID, m.Title, mgErr)
		}
	}

	for _, actorID := range m.ActorId {
		actor, err := actorRepo.FindById(actorID)
		if err != nil || actor == nil {
			continue
		}
		if _, maErr := r.DB.Exec(maQuery, movieID, actorID); maErr != nil {
			return 0, fmt.Errorf("failed to link actor %d to MOVIES %q: %w", actorID, m.Title, maErr)
		}
	}
	return movieID, nil
}

func (r *MovieRepository) Patch() {
}

func (r *MovieRepository) Delete(MovieId int, force bool) (int64, error) {
	if force {
		m, payloadErr := r.GetById(MovieId)
		if payloadErr != nil {
			return 0, payloadErr
		}
		mgQuery := `DELETE FROM movie_genre WHERE genreId = ?`
		maQuery := `DELETE FROM movie_actor WHERE actorId = ?`
		genresId, _ := r.GetMovie_genre(MovieId)
		for _, genreID := range genresId {
			if _, mgErr := r.DB.Exec(mgQuery, genreID); mgErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %s: %w", genreID, m.Title, mgErr)
			}
		}
		actorsId, _ := r.GetMovie_actor(MovieId)
		for _, actorID := range actorsId {
			if _, maErr := r.DB.Exec(maQuery, actorID); maErr != nil {
				return 0, fmt.Errorf("failed to remove the link bw genre %d and MOVIES %s: %w", actorID, m.Title, maErr)
			}
		}
	}
	rows, deleteErr := r.DB.Exec("DELETE FROM MOVIES WHERE id=?", MovieId)
	if deleteErr != nil {
		return 0, deleteErr
	}
	affectedR, _ := rows.RowsAffected()
	return affectedR, nil
}

// this method returns the genreIds of provided movieId
func (r *MovieRepository) GetMovie_genre(MovieId int) ([]int, int) {
	count := 0
	genresId := []int{}
	mgQuery := `SELECT genreId FROM movie_genre WHERE movieId = ?`
	rows, err := r.DB.Query(mgQuery, MovieId)
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
func (r *MovieRepository) GetMovie_actor(MovieId int) ([]int, int) {
	count := 0
	actorsId := []int{}
	maQuery := `SELECT actorId FROM movie_actor WHERE movieId = ?`

	rows, err := r.DB.Query(maQuery, MovieId)
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
func (r *MovieRepository) GetGenresForMovie(MovieId int) ([]string, error) {
	movieGenres := []string{}
	genreRepo := NewGenreRepository(r.DB)
	genresId, _ := r.GetMovie_genre(MovieId)

	for _, genreId := range genresId {
		genre, err := genreRepo.GetGenreByID(genreId)
		if err != nil {
			continue
		}
		movieGenres = append(movieGenres, genre.Name)
	}
	return movieGenres, nil
}

// this method returnes the names of actors of a provided movieId by using the GetMovie_actor function
func (r *MovieRepository) GetActorsForMovie(MovieId int) ([]string, error) {
	movieActors := []string{}
	actorRepo := NewActorRepository(r.DB)
	actorsId, _ := r.GetMovie_actor(MovieId)

	for _, actorid := range actorsId {
		actor, err := actorRepo.FindById(actorid)
		if err != nil {
			continue
		}
		movieActors = append(movieActors, actor.Name)
	}

	return movieActors, nil
}
func formatDurations(duration uint16) string {
	return ""
}
