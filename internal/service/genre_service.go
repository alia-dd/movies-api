package service

import (
	"fmt"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"movies-api/internal/repository"
	"strings"
)

type GenreService struct {
	repo *repository.GenreRepository
}

func NewGenreService(repo *repository.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func (s *GenreService) CreateGenre(genre *models.Genre) error {
	if genre.Name == "" {
		return errors.ErrInvalidInput
	}
	genre.Name = strings.TrimSpace(genre.Name)

	return s.repo.CreateGenre(genre)
}

func (s *GenreService) GetAllGenres(page, limit int) ([]models.Genre, int, error) {
	if page <= 0 || limit <= 0 {
		return nil, 0, errors.ErrInvalidInput
	}

	return s.repo.GetAllGenresWithPagination(page, limit)
}

func (s *GenreService) GetGenreByID(id int) (*models.Genre, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetGenreByID(id)
}

func (s *GenreService) SearchGenreByName(searchTerm string, page, limit int) ([]models.Genre, int, error) {
	if searchTerm == "" {
		return nil, 0, errors.ErrInvalidInput
	}
	if page <= 0 || limit <= 0 {
		return nil, 0, errors.ErrInvalidInput
	}

	searchTerm = strings.TrimSpace(searchTerm)
	return s.repo.SearchGenreByName(searchTerm, page, limit)
}

func (s *GenreService) UpdateGenre(id int, name string) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}
	if name == "" {
		return errors.ErrInvalidInput
	}

	name = strings.TrimSpace(name)
	return s.repo.UpdateGenre(id, name)
}

func (s *GenreService) DeleteGenreByID(id int, force bool) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if !force {
		count, err := s.repo.CountMoviesByGenre(id)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete genre because it has %d associated movies. Use ?force=true to force delete", count)
		}
	}

	return s.repo.DeleteGenreWithAssociations(id)
}

func (s *GenreService) DeleteGenreByName(name string, force bool) error {
	if name == "" {
		return errors.ErrInvalidInput
	}

	name = strings.TrimSpace(name)

	genre, err := s.repo.GetGenreByName(name)
	if err != nil {
		return err
	}

	if !force {
		count, err := s.repo.CountMoviesByGenre(genre.Id)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete genre '%s' because it has %d associated movies. Use ?force=true to force delete", name, count)
		}
	}

	return s.repo.DeleteGenreByNameWithAssociations(name)
}
