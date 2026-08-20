package repository

import (
	"database/sql"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"strings"
	"time"
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
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.ErrDuplicateKey
		}
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
		return errors.ErrNotFound
	}
	actor.UpdatedAt = now
	return nil
}

func (r *ActorsRepository) FindById(id int) (*models.Actor, error) {
	query := `
	SELECT id,name,birthdate FROM actors WHERE id = ?`
	actor := &models.Actor{}

	err := r.db.QueryRow(query, id).Scan(
		&actor.Id,
		&actor.Name,
		&actor.BirthDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return actor, nil
}

func (r *ActorsRepository) FindByName(name string) (*models.Actor, error) {
	query := `
	SELECT id,name,birthdate FROM actors WHERE name = ?`
	actor := &models.Actor{}

	err := r.db.QueryRow(query, name).Scan(
		&actor.Id,
		&actor.Name,
		&actor.BirthDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return actor, nil
}

func (r *ActorsRepository) GetAllActors() ([]models.Actor, error) {
	query := `
	SELECT id, name, birthdate FROM actors`

	actors := []models.Actor{}
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, err
}

func (r *ActorsRepository) DeleteActorsById(id int) error {
	query := `
	DELETE FROM actors WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsDeleted == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *ActorsRepository) DeleteActorsByName(name string) error {
	query := `
	DELETE FROM actors WHERE name = ?`

	result, err := r.db.Exec(query, name)
	if err != nil {
		return err
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsDeleted == 0 {
		return errors.ErrNotFound
	}
	return nil
}
