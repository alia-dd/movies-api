package service

import (
	"log"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type MovieService struct {
	repo *repository.DatabaseConnection
}

// meshan waxaa uguwaxdaa repoga haa waxan dhan waakabodi karate
func NewMovieService(repo *repository.DatabaseConnection) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) GetMovie() ([]models.Movies, error) {
	payload, err := s.repo.Get()
	return payload, err
}

func (s *MovieService) GetMovieById(id int) (*models.Movies, error) {
	payload, err := s.repo.GetById(id)
	return payload, err
}
func (s *MovieService) GetMovieByTitle(actorName string) (*models.Movies, error) {
	payload, err := s.repo.GetByTitle(actorName)
	return payload, err
}

func (s *MovieService) DeleteMovie(id int) error {
	affectedR, err := s.repo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete movie with id = %d \n Error: %v", id, err)
		return err
	}
	log.Printf("Successfully deleted Number of Row Affected %d", affectedR)

	return nil
}
