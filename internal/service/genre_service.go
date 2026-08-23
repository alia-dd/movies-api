package service

import (
	"context"
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

func (s *GenreService) CreateGenre(cx context.Context, genre *models.Genre) error {
	if genre.Name == "" {
		return errors.ErrInvalidInput
	}
	genre.Name = strings.TrimSpace(genre.Name)

	return s.repo.CreateGenre(cx, genre)
}

func (s *GenreService) GetAllGenres(cx context.Context, page, limit int) ([]models.Genre, int, error) {
	if page <= 0 || limit <= 0 {
		return nil, 0, errors.ErrInvalidInput
	}

	return s.repo.GetAllGenresWithPagination(cx, page, limit)
}

func (s *GenreService) GetGenreByID(cx context.Context, id int) (*models.Genre, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetGenreByID(cx, id)
}

func (s *GenreService) SearchGenreByName(cx context.Context, searchTerm string, page, limit int) ([]models.Genre, int, error) {
	if searchTerm == "" {
		return nil, 0, errors.ErrInvalidInput
	}
	if page <= 0 || limit <= 0 {
		return nil, 0, errors.ErrInvalidInput
	}

	searchTerm = strings.TrimSpace(searchTerm)
	return s.repo.SearchGenreByName(cx, searchTerm, page, limit)
}

func (s *GenreService) UpdateGenre(cx context.Context, id int, name string) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}
	if name == "" {
		return errors.ErrInvalidInput
	}

	name = strings.TrimSpace(name)
	return s.repo.UpdateGenre(cx, id, name)
}

func (s *GenreService) DeleteGenreByID(cx context.Context, id int, force bool) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if !force {
		count, err := s.repo.CountMoviesByGenre(cx, id)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete genre because it has %d associated movies. Use ?force=true to force delete", count)
		}
	}

	return s.repo.DeleteGenreWithAssociations(cx, id)
}

func (s *GenreService) DeleteGenreByName(cx context.Context, name string, force bool) error {
	if name == "" {
		return errors.ErrInvalidInput
	}

	name = strings.TrimSpace(name)

	genre, err := s.repo.GetGenreByName(cx, name)
	if err != nil {
		return err
	}

	if !force {
		count, err := s.repo.CountMoviesByGenre(cx, genre.Id)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete genre '%s' because it has %d associated movies. Use ?force=true to force delete", name, count)
		}
	}

	return s.repo.DeleteGenreByNameWithAssociations(cx, name)
}
