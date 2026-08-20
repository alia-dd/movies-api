package service

import (
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"movies-api/internal/errors"

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

func (s *ActorService) CreateActor(actor *models.Actor) error {
	if actor.Name == "" || actor.BirthDate == "" {
		return errors.ErrInvalidInput
	}
	_, err := time.Parse("2006-01-02", actor.BirthDate)
	if err != nil {
		return errors.ErrInvalidInput
	}
	return s.repo.CreateActor(actor)

}

func (s *ActorService) UpdateActor(actor *models.Actor) error {
	if actor.Name == "" || actor.BirthDate == "" {
		return errors.ErrInvalidInput
	}
	_, err := time.Parse("2006-01-02", actor.BirthDate)
	if err != nil {
		return errors.ErrInvalidInput
	}
	return s.repo.Update(actor)
}

func (s *ActorService) FindById(id int) (*models.Actor, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}
	return s.repo.FindById(id)
}
func (s *ActorService) FindByName(name string) (*models.Actor, error) {
	if name == "" {
		return nil, errors.ErrInvalidInput
	}
	return s.repo.FindByName(name)
}

func (s *ActorService) GetAllActors(page, limit int) ([]models.Actor, int, error) {
	return s.repo.GetAllActors(page, limit)
}
func (s *ActorService) SearchActorByName(search string, page, limit int) ([]models.Actor, int, error) {
	return s.repo.SearchActorByName(search, page, limit)
}
func (s *ActorService) DeleteActorsById(id int, force bool) error {
	return s.repo.DeleteActorsById(id, force)
}
func (s *ActorService) DeleteActorsByName(name string, force bool) error {
	if name == "" {
		return errors.ErrInvalidInput
	}
	return s.repo.DeleteActorsByName(name, force)
}
