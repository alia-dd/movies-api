package service

import (
	"context"
	"fmt"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"strings"
	"time"
)

type ActorService struct {
	repo *repository.ActorsRepository
}

func NewActorService(repo *repository.ActorsRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

func (s *ActorService) CreateActor(ctx context.Context, actor *models.Actor) error {
	if actor.Name == "" || actor.BirthDate == "" {
		return errors.ErrInvalidInput
	}
	_, err := time.Parse("2006-01-02", actor.BirthDate)
	if err != nil {
		return errors.ErrInvalidInput
	}
	return s.repo.CreateActor(ctx, actor)

}

func (s *ActorService) UpdateActor(ctx context.Context, id int, actor *models.UpdateActor) error {
	if id <= 0 {
		fmt.Println("id problem")
		return errors.ErrInvalidInput
	}
	if actor.Name == nil && actor.BirthDate == nil {
		fmt.Println("nil problem")
		return errors.ErrInvalidInput
	}
	if actor.Name != nil && strings.TrimSpace(*actor.Name) == "" {
		fmt.Println("trim problem")
		return errors.ErrInvalidInput
	}
	if actor.BirthDate != nil {
		_, err := time.Parse("2006-01-02", *actor.BirthDate)
		if err != nil {
			fmt.Println("time problem")
			return errors.ErrInvalidInput
		}
	}

	return s.repo.Update(ctx, id, actor)
}

func (s *ActorService) FindById(ctx context.Context, id int) (*models.Actor, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}
	return s.repo.FindById(ctx, id)
}
func (s *ActorService) FindByName(ctx context.Context, name string) (*models.Actor, error) {
	if name == "" {
		return nil, errors.ErrInvalidInput
	}
	return s.repo.FindByName(ctx, name)
}

func (s *ActorService) GetAllActors(ctx context.Context, page, limit int) ([]models.Actor, int, error) {
	return s.repo.GetAllActors(ctx, page, limit)
}

func (s *ActorService) DeleteActorsById(ctx context.Context, id int, force bool) error {
	return s.repo.DeleteActorsById(ctx, id, force)
}
func (s *ActorService) DeleteActorsByName(ctx context.Context, name string, force bool) error {
	if name == "" {
		return errors.ErrInvalidInput
	}
	return s.repo.DeleteActorsByName(ctx, name, force)
}
