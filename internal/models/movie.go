package models

import (
	"time"
)

type Movies struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseYear int       `json:"releaseyear"`
	Duration    int       `json:"duration"`
	GenreId     []int     `json:"genreid"`
	ActorId     []int     `json:"actorid"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
}
