package repository

import (
	"database/sql"
	"errors"
	"moviesApi/internal/models"
	"time"
)

var (
	ErrNotFound     = errors.New("record Not Found")
	ErrDuplicateKey = errors.New("duplicate key violaion")
	ErrInvalidInput = errors.New("invalid input")
)

type ActorsRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorsRepository {
	return &ActorsRepository{
		db: db,
	}
}

func (r *ActorsRepository) CreateActor(actor *models.Actor) error {

	query := `
	INSERT INTO actors(name, birthdate)
	VALUES(?,?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, actor.Name, actor.BirthDate)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	actor.Id = int(id)
	actor.CreatedAt = now
	actor.UpdatedAt = now
	return nil
}

func (r *ActorsRepository) Update(actor *models.Actor) error {
	query := `
	UPDATE actors
	SET name = ?, birthdate = ?, updated_at = ?
	WHERE id = ?
	`
	now := time.Now()
	result, err := r.db.Exec(query, actor.Name, actor.BirthDate, now, actor.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	actor.UpdatedAt = now
	return nil
}

func (r *ActorsRepository) FindById(id int) (*models.Actor, error) {
	query := `
	SELECT id,name,birthdate
	FROM actors
	WHERE id = ?
	`
	actor := &models.Actor{}

	err := r.db.QueryRow(query, id).Scan(
		&actor.Id,
		&actor.Name,
		&actor.BirthDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return actor, nil
}

// func (r *ActorsRepository) FilterByName(name string)(models.Actor, err){

// }
