package models

import "time"

type Actor struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birthDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateActor struct {
	Name      *string `json: "name, omitempty"`
	BirthDate *string `"json: birthdate, omitempty"`
}
