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
	ReleaseYear int    `json:"release_date"`

	Duration         int    `json:"duration"`
	Overview         string `json:"overview"`
	OriginalLanguage string `json:"original_language"`

	GenreId []int `json:"genre_ids"`
	ActorId []int `json:"Actor_ids"`

	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
