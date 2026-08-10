package service

import (
	"fmt"
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
	if f.Genre != "" {
		genreID, err := strconv.Atoi(f.Genre)
		if err != nil || genreID < 0 {
			log.Printf("Invalid Genre ID: Genre must be zero or a positive integer")
			f.Genre = ""
		}
	}
	if f.Year != "" {
		year, err := strconv.Atoi(f.Year)
		if err != nil || year < 1888 || year > time.Now().Year() {
			log.Printf("Invalid Yead")
			f.Year = ""
		}
	}
	if f.Actor != "" {
		ActorID, err := strconv.Atoi(f.Actor)
		if err != nil || ActorID < 0 {
			log.Printf("Invalid Genre ID: Genre must be zero or a positive integer")
			f.Actor = ""
		}
	}
	if f.Page != "" {
		page, err := strconv.Atoi(f.Page)
		if err != nil || page <= 0 {
			log.Printf("Invalid page number")
			f.Page = ""
		}
	}
	if f.Size != "" {
		size, err := strconv.Atoi(f.Size)
		if err != nil || size <= 0 {
			log.Printf("Invalid size number")
			f.Size = ""
		}
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

func (s *MovieService) CreateMovie(movie models.Movies) error {
	if movie.Title == "" || movie.ReleaseYear < 1888 || movie.ReleaseYear > time.Now().Year() {
		return fmt.Errorf("Invalid input")
	}
	fmt.Println("herre  CreateMovie")
	err := s.repo.Post(movie)
	return err
}
func (s *MovieService) DeleteMovie(id string, force string) error {
	forcebool, forceErr := strconv.ParseBool(force)
	if forceErr != nil {
		log.Printf("force must be a boolean")
		return forceErr
	}
	var sentencedMovie *models.Movies
	var sentencederr error
	mId, idErr := strconv.Atoi(id)
	if idErr != nil {
		sentencedMovie, sentencederr = s.GetMovieByTitle(id)
	} else {
		sentencedMovie, sentencederr = s.GetMovieById(mId)
	}
	if sentencederr != nil {
		return sentencederr
	}
	gid := len(sentencedMovie.GenreId)
	aId := len(sentencedMovie.ActorId)
	if !forcebool && gid > 0 {
		return fmt.Errorf("Cannot delete movie %s because it has %d associated genre", id, gid)
	}
	if !forcebool && aId > 0 {
		return fmt.Errorf("Cannot delete movie %s because it has %d associated actors", id, aId)
	}
	affectedR, err := s.repo.Delete(mId, forcebool)
	if err != nil {
		log.Printf("Failed to delete movie with id = %d \n Error: %v", id, err)
		return err
	}
	log.Printf("Successfully deleted Number of Row Affected %d", affectedR)

	return nil
}
