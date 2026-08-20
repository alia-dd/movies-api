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
	return s.repo.FindById(id)
}
func (s *ActorService) FindByName(name string) (*models.Actor, error) {
	if name == "" {
		return nil, errors.ErrInvalidInput
	}
	return s.repo.FindByName(name)
}

func (s *ActorService) GetAllActors() ([]models.Actor, error) {
	return s.repo.GetAllActors()
}
func (s *ActorService) DeleteActorsById(id int) error {
	return s.repo.DeleteActorsById(id)
}
func (s *ActorService) DeleteActorsByName(name string) error {
	if name == "" {
		return errors.ErrInvalidInput
	}
	return s.repo.DeleteActorsByName(name)
}
