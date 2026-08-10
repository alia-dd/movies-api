package repository

import (
	"database/sql"
	"moviesApi/internal/models"
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

func (r *GenreRepository) CreateGenre(genre *models.Genre) error {
	query := `INSERT INTO genres (name) VALUES(?)`
	result, err := r.db.Exec(query, genre.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateKey
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	genre.Id = int(id)

	err = r.db.QueryRow(`SELECT created_at, updated_at FROM genres WHERE id = ?`, id).
		Scan(&genre.CreatedAt, &genre.UpdatedAt)
	return err

	/*
		line 42 : why
		if not used the fetched data in stdout create and update times show default time
		{"id":3,"name":"Drama","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}
		to update our data
		we fetch the row from db by id
		using Scan() we assign the data into corresponding place
		order of variables important
	*/
}

func (r *GenreRepository) GetAllGenres() ([]models.Genre, error) {
	query := `SELECT * FROM genres`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.Id, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)
		if err !=nil {
			return nil,err
		}
		genres =append(genres,genre)
	}
	return genres,rows.Err()
}

func (r *GenreRepository) GetGenreByID(id int) (*models.Genre,error) {
	query := `SELECT * FROM genres WHERE id = ?`
	var genre models.Genre
	err :=r.db.QueryRow(query,id).Scan(&genre.Id,&genre.Name,&genre.CreatedAt,&genre.UpdatedAt)
	if err == sql.ErrNoRows{
		return nil,ErrNotFound
	}
	if err !=nil {
		return nil,err
	}
	return &genre,nil
}

func (r *GenreRepository) UpdateGenre(id int, name string) error {
	query := `UPDATE genres SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result,err := r.db.Exec(query,name,id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateKey
		}
		return err
	}

	rows,err := result.RowsAffected() //affected number of rows
	if err !=nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}


func(r *GenreRepository) DeleteGenreByID(id int)error {
	query := `DELETE FROM genres WHERE id = ?`
	result,err := r.db.Exec(query,id)
	if err != nil {
		return err
	}
	rows,err := result.RowsAffected()
	if err != nil {
		return err
	} 
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func(r *GenreRepository) DeleteGenreByName(name string) error {
	query := `DELETE FROM genres WHERE name = ?`
	result,err := r.db.Exec(query,name)
	if err != nil {
		return err
	}
	rows,err := result.RowsAffected()
	if err != nil {
		return err
	} 
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}