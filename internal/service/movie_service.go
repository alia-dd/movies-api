package service

import (
	"log"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"strconv"
	"time"
)

type MovieService struct {
	repo *repository.DatabaseConnection
}

// meshan waxaa uguwaxdaa repoga haa waxan dhan waakabodi karate
func NewMovieService(repo *repository.DatabaseConnection) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) GetMovie(f models.Filter) ([]models.Movies, error) {
	genreID, err := strconv.Atoi(f.Genre)
	if err != nil || genreID < 0 {
		log.Printf("Invalid Genre ID: Genre must be zero or a positive integer")
		f.Genre = ""
	}
	year, err := strconv.Atoi(f.Year)
	if err != nil || year < 1888 || year > time.Now().Year() {
		log.Printf("Invalid Yead")
		f.Year = ""
	}
	ActorID, err := strconv.Atoi(f.Actor)
	if err != nil || ActorID < 0 {
		log.Printf("Invalid Genre ID: Genre must be zero or a positive integer")
		f.Actor = ""
	}

	page, err := strconv.Atoi(f.Page)
	if err != nil || page <= 0 {
		log.Printf("Invalid page number")
		f.Page = ""
	}
	size, err := strconv.Atoi(f.Size)
	if err != nil || size <= 0 {
		log.Printf("Invalid size number")
		f.Size = ""
	}

	payload, err := s.repo.Get(f)
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

func (s *MovieService) DeleteMovie(id int, force string) error {
	forcebool, forceErr := strconv.ParseBool(force)
	if forceErr != nil {
		log.Printf("force must be a boolean")
		return forceErr
	}
	affectedR, err := s.repo.Delete(id, forcebool)
	if err != nil {
		log.Printf("Failed to delete movie with id = %d \n Error: %v", id, err)
		return err
	}
	log.Printf("Successfully deleted Number of Row Affected %d", affectedR)

	return nil
}
