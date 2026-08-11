package service

import (
	"moviesApi/internal/models"
	"moviesApi/internal/repository"
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
		return repository.ErrInvalidInput
	}
	return s.repo.CreateGenre(genre)
}

func (s *GenreService) GetAllGenres(page,limit int) ([]models.Genre, int, error) {
	return s.repo.GetAllGenres(page,limit)
}

func (s *GenreService) GetGenreByID(id int) (*models.Genre, error) {
	if id <= 0 {
		return nil, repository.ErrInvalidInput
	}
	return s.repo.GetGenreByID(id)
}

func (s *GenreService) UpdateGenre(id int, nameInput string) error {
	name := strings.TrimSpace(nameInput)
	if id <= 0 || name == "" {
		return repository.ErrInvalidInput
	}
	return s.repo.UpdateGenre(id, name)
}

func (s *GenreService) DeleteGenreByID(id int) error {
	if id <= 0 {
		return repository.ErrInvalidInput
	}
	return s.repo.DeleteGenreByID(id)
}

func (s *GenreService) DeleteGenreByName(nameInput string) error {
	name := strings.TrimSpace(nameInput)
	if name == "" {
		return repository.ErrInvalidInput
	}
	return s.repo.DeleteGenreByName(name)
}
