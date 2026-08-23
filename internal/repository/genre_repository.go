package repository

import (
	"context"
	"database/sql"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"strings"
)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

func (r *GenreRepository) CreateGenre(cx context.Context, genre *models.Genre) error {
	query := `INSERT INTO GENRES (name) VALUES(?)`
	result, err := r.db.ExecContext(cx, query, genre.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.ErrDuplicateKey
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	genre.Id = int(id)

	err = r.db.QueryRowContext(cx, `SELECT created_at, updated_at FROM GENRES WHERE id = ?`, id).
		Scan(&genre.CreatedAt, &genre.UpdatedAt)
	return err
}

func (r *GenreRepository) GetAllGenresWithPagination(cx context.Context, page, limit int) ([]models.Genre, int, error) {
	countQuery := `SELECT COUNT(*) FROM GENRES`
	var total int
	err := r.db.QueryRowContext(cx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, name, created_at, updated_at FROM GENRES LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(cx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.Id, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		genres = append(genres, genre)
	}
	return genres, total, rows.Err()
}

func (r *GenreRepository) GetGenreByID(cx context.Context, id int) (*models.Genre, error) {
	query := `SELECT id, name, created_at, updated_at FROM GENRES WHERE id = ?`
	var genre models.Genre
	err := r.db.QueryRowContext(cx, query, id).Scan(&genre.Id, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

func (r *GenreRepository) GetGenreByName(cx context.Context, name string) (*models.Genre, error) {
	query := `SELECT id, name, created_at, updated_at FROM GENRES WHERE name = ?`
	var genre models.Genre
	err := r.db.QueryRowContext(cx, query, name).Scan(&genre.Id, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	return &genre, err
}

func (r *GenreRepository) SearchGenreByName(cx context.Context, searchTerm string, page, limit int) ([]models.Genre, int, error) {
	countQuery := `SELECT COUNT(*) FROM GENRES WHERE name LIKE ?`
	var total int
	err := r.db.QueryRowContext(cx, countQuery, "%"+searchTerm+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, name, created_at, updated_at FROM GENRES WHERE name LIKE ? LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(cx, query, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.Id, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		genres = append(genres, genre)
	}
	return genres, total, rows.Err()
}

func (r *GenreRepository) CountMoviesByGenre(cx context.Context, genreId int) (int, error) {
	query := `SELECT COUNT(*) FROM movie_genre WHERE genreId = ?`
	var count int
	err := r.db.QueryRowContext(cx, query, genreId).Scan(&count)
	return count, err
}

func (r *GenreRepository) UpdateGenre(cx context.Context, id int, name string) error {
	query := `UPDATE GENRES SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := r.db.ExecContext(cx, query, name, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.ErrDuplicateKey
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *GenreRepository) DeleteGenre(cx context.Context, id int) error {
	query := `DELETE FROM GENRES WHERE id = ?`
	result, err := r.db.ExecContext(cx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *GenreRepository) DeleteGenreWithAssociations(cx context.Context, id int) error {
	query1 := `DELETE FROM movie_genre WHERE genreId = ?`
	_, err := r.db.ExecContext(cx, query1, id)
	if err != nil {
		return err
	}

	query2 := `DELETE FROM GENRES WHERE id = ?`
	result, err := r.db.ExecContext(cx, query2, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *GenreRepository) DeleteGenreByName(cx context.Context, name string) error {
	query := `DELETE FROM GENRES WHERE name = ?`
	result, err := r.db.ExecContext(cx, query, name)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *GenreRepository) DeleteGenreByNameWithAssociations(cx context.Context, name string) error {
	var genreId int
	err := r.db.QueryRowContext(cx, `SELECT id FROM GENRES WHERE name = ?`, name).Scan(&genreId)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(cx, `DELETE FROM movie_genre WHERE genreId = ?`, genreId)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(cx, `DELETE FROM GENRES WHERE id = ?`, genreId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}
