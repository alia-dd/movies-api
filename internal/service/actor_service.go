package service

import (
	"moviesApi/internal/models"
	"moviesApi/internal/repository"
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
		return repository.ErrInvalidInput
	}

	return s.repo.CreateActor(actor)

}

func (s *ActorService) UpdateActor(actor *models.Actor) error {
	if actor.Name == "" {
		return repository.ErrInvalidInput
	}
	return s.repo.Update(actor)
}

func (s *ActorService) FindById(id int) (*models.Actor, error) {
	return s.repo.FindById(id)
}
func (s *ActorService) FindByName(name string) (*models.Actor, error) {
	if name == "" {
		return nil, repository.ErrInvalidInput
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
	return s.repo.DeleteActorsByName(name)
}
