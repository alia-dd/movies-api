package repository

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *ActorsRepository) CreateActor(ctx context.Context, actor *models.Actor) error {
	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		return errors.ErrTransactionStart
	}
	defer tx.Rollback() // if there is error revert changes to the db back to before the changes

	query := `INSERT INTO actors(name, birthdate) VALUES(?,?)`

	now := time.Now()
	result, err := tx.ExecContext(ctx, query, actor.Name, actor.BirthDate)
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

func (r *ActorsRepository) Update(ctx context.Context, actor *models.Actor) error {
	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		return errors.ErrTransactionStart
	}
	defer tx.Rollback() // if there is error revert changes to the db back to before the changes
	query := `UPDATE actors	SET name = ?, birthdate = ?, updated_at = ?	WHERE id = ?`
	now := time.Now()
	result, err := tx.ExecContext(ctx, query, actor.Name, actor.BirthDate, now, actor.Id)
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

func (r *ActorsRepository) FindById(ctx context.Context, id int) (*models.Actor, error) {
	query := `SELECT id,name,birthdate FROM actors WHERE id = ?`
	actor := &models.Actor{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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

func (r *ActorsRepository) FindByName(ctx context.Context, name string) (*models.Actor, error) {
	query := `SELECT id,name,birthdate FROM actors WHERE name = ?`
	actor := &models.Actor{}

	err := r.db.QueryRowContext(ctx, query, name).Scan(
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

func (r *ActorsRepository) GetAllActors(ctx context.Context, page, limit int) ([]models.Actor, int, error) {
	offset := (page - 1) * limit
	query := `SELECT * FROM actors LIMIT ? OFFSET ?`
	var total int
	count := `SELECT COUNT(*) FROM actors`
	r.db.QueryRowContext(ctx, count).Scan(&total)
	rows, err := r.db.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	actors := []models.Actor{}
	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return nil, total, err
		}
		actors = append(actors, actor)
	}

	return actors, total, err
}

func (r *ActorsRepository) DeleteActorsById(ctx context.Context, id int, force bool) error {
	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		return errors.ErrTransactionStart
	}
	defer tx.Rollback() // if there is error revert changes to the db back to before the changes
	// checking if the actor exists
	var name string
	err := tx.QueryRowContext(ctx, "SELECT name FROM actors WHERE id = ?", id).Scan(&name)
	if err != nil {
		return errors.ErrNotFound
	}
	// checking if the actor is related to any movies
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM movie_actors WHERE actor_id = ?", id).Scan(&count)

	if err != nil {
		return err
	}
	// if force deletion isn't requested and actor has relation to movies
	if count > 0 && !force {
		return fmt.Errorf("can't delete actor %s, because it has %d related movies", name, count)
	}
	// removing the relation first
	if force {
		_, err := tx.ExecContext(ctx, `DELETE FROM movie_actors WHERE actor_id = ?`, id)
		if err != nil {
			return err
		}
	}
	// delete the actor
	query := `
	DELETE FROM actors WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
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

// not needed because if we have multiple actors with the same name, it doesn't know which one to delete
func (r *ActorsRepository) DeleteActorsByName(ctx context.Context, name string, force bool) error {
	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		return errors.ErrTransactionStart
	}
	defer tx.Rollback() // if there is error revert changes to the db back to before the changes
	// checking if the actor exists by exracting the name using Id
	var id int
	err := tx.QueryRowContext(ctx, "SELECT id FROM actors WHERE name = ?", name).Scan(&id)
	if err != nil {
		return errors.ErrNotFound
	}
	// checking if the actor is related to any movies
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM movie_actors WHERE actor_id = ?", id).Scan(&count)

	if err != nil {
		return err
	}
	// if force deletion isn't requested and actor has relation to movies
	if count > 0 && !force {
		return fmt.Errorf("can't delete actor %s, because it has %d related movies", name, count)
	}
	// removing the relation first
	if force {
		_, err := tx.Exec(`DELETE FROM movie_actors WHERE actor_id = ?`, id)
		if err != nil {
			return err
		}
	}
	// delete the actor
	query := `DELETE FROM actors WHERE id = ?`

	result, err := tx.ExecContext(ctx, query, id)
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
