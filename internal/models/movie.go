package models

import (
	"time"
)

type Filter struct {
	Genre string
	Year  string
	Actor string
	Page  string
	Size  string
}

type Movies struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseYear string `json:"release_date"`

	Duration         uint16 `json:"duration"`
	Overview         string `json:"overview"`
	OriginalLanguage string `json:"original_language"`

	GenreId []int `json:"genre_ids"`
	ActorId []int `json:"Actor_ids"`

	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type MovieUpdate struct {
	Title       *string `json:"title,omitempty"`
	ReleaseYear *int    `json:"releaseYear,omitempty"`

	Duration         *int    `json:"duration,omitempty"`
	Overview         *string `json:"overview"`
	OriginalLanguage *string `json:"original_language"`

	AddActorIDs    []int `json:"addActorIds,omitempty"`
	RemoveActorIDs []int `json:"removeActorIds,omitempty"`
	AddGenreIDs    []int `json:"addGenreIds,omitempty"`
	RemoveGenreIDs []int `json:"removeGenreIds,omitempty"`
}
type MoviesDisplay struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseYear string `json:"release_date"`

	Duration         string `json:"duration"`
	Overview         string `json:"overview"`
	OriginalLanguage string `json:"original_language"`

	Genres []string `json:"genre_ids"`
	Actors []string `json:"Actor_ids"`

	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
