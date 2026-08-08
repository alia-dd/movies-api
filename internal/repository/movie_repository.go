package repository

import (
	"database/sql"
	"movies-api/internal/models"
)

type DatabaseConnection struct {
	DB *sql.DB
}

// meshan waxaa galaayo sql strings data
func NewdbConnection(db *sql.DB) *DatabaseConnection {
	return &DatabaseConnection{DB: db}
}
func (c *DatabaseConnection) Get() ([]models.Movies, error) {
	rows, rowErr := c.DB.Query("SELECT * FROM movie")
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	movies := []models.Movies{}

	for rows.Next() {
		movie := models.Movies{}
		err := rows.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	// if rowErr == sql.ErrNoRows {
	// 	return nil, rowErr
	// }

	return movies, nil
}
func (c *DatabaseConnection) GetById(id int) (*models.Movies, error) {
	row := c.DB.QueryRow("SELECT FROM movie WHERE id=? ", id)
	movie := models.Movies{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (c *DatabaseConnection) GetByTitle(actorName string) (*models.Movies, error) {
	row := c.DB.QueryRow("SELECT FROM movie WHERE title=? ", actorName)
	movie := models.Movies{}
	err := row.Scan(&movie.Id, &movie.Title, &movie.ReleaseYear, &movie.Duration, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (c *DatabaseConnection) Post() {
}

func (c *DatabaseConnection) Patch() {
}

func (c *DatabaseConnection) Delete(id int) (int64, error) {
	rows, deleteErr := c.DB.Exec("DELETE FROM movie WHERE id=?", id)
	if deleteErr != nil {
		return 0, deleteErr
	}
	affectedR, _ := rows.RowsAffected()
	return affectedR, nil
}
