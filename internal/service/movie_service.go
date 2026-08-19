package service

import (
	"context"
	"fmt"
	"log"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"strconv"
	"time"
)

type MovieService struct {
	repo *repository.MovieRepository
}

// meshan waxaa uguwaxdaa repoga haa waxan dhan waakabodi karate
func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) GetMovie(cx context.Context, f models.Filter) ([]models.MoviesDisplay, error) {
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
	payload, err := s.repo.Get(cx, &f)
	return payload, err
}

func (s *MovieService) GetMovieById(cx context.Context, id int) (*models.MoviesDisplay, error) {
	payload, err := s.repo.GetById(cx, id)
	return payload, err
}
func (s *MovieService) GetMovieByTitle(cx context.Context, title string) (*models.MoviesDisplay, error) {
	payload, err := s.repo.GetByTitle(cx, title)
	return payload, err
}
func (s *MovieService) SearchMovieByTitle(cx context.Context, title string) ([]models.MoviesDisplay, error) {
	payload, err := s.repo.SearchByTitle(cx, title)
	return payload, err
}

func (s *MovieService) CreateMovie(cx context.Context, movie models.Movies) (*models.MoviesDisplay, error) {
	year, err := strconv.Atoi(movie.ReleaseYear)
	if err != nil || movie.Title == "" || year < 1888 || year > time.Now().Year() {
		return nil, fmt.Errorf("Invalid input")
	}
	movieId, postErr := s.repo.Post(cx, movie)
	if postErr != nil {
		return nil, postErr
	}
	return s.repo.GetById(cx, int(movieId))
}

func (s *MovieService) UpdateMovie(cx context.Context, id int, m models.MovieUpdate) (*models.MoviesDisplay, error) {
	// handle check here
	if _, err := s.repo.GetById(cx, id); err != nil {
		return nil, err
	}
	if patchErr := s.repo.Patch(cx, id, m); patchErr != nil {
		return nil, patchErr
	}
	return s.repo.GetById(cx, id)
}

func (s *MovieService) DeleteMovie(cx context.Context, id string, force string) error {
	forcebool, forceErr := strconv.ParseBool(force)
	if forceErr != nil {
		return fmt.Errorf("force must be a boolean %w", forceErr)
	}
	var sentencederr error
	fmt.Println("here")
	mId, idErr := strconv.Atoi(id)
	if idErr != nil {
		var sentencedMovie *models.MoviesDisplay
		sentencedMovie, sentencederr = s.GetMovieByTitle(cx, id)
		mId = sentencedMovie.Id
	}
	if sentencederr != nil {
		return sentencederr
	}

	_, gCount := s.repo.GetMovie_genre(cx, mId)
	_, aCount := s.repo.GetMovie_actor(cx, mId)
	fmt.Println(gCount, aCount)
	if !forcebool && gCount > 0 {
		return fmt.Errorf("Cannot delete movie %s because it has %d associated genre", id, gCount)
	}
	if !forcebool && aCount > 0 {
		return fmt.Errorf("Cannot delete movie %s because it has %d associated actors", id, aCount)
	}
	affectedR, err := s.repo.Delete(cx, mId, forcebool)
	if err != nil {
		return fmt.Errorf("Failed to delete movie with id = %v \n Error: %v", id, err)
	}
	log.Printf("Successfully deleted Number of Row Affected %d", affectedR)

	return nil
}
