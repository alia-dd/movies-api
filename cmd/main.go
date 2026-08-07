package main

import (
	"fmt"
	"log"
	"moviesApi/internal/database"
	"moviesApi/internal/models"
	"moviesApi/internal/repository"
	// "time"
)

func main() {
	db, err := database.OpenDB("movies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = database.NewTable(db)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println("table created")

	repo := repository.NewActorRepository(db)

	actor := &models.Actor{
		Name:      "whaaaat Hanks",
		BirthDate: "1956-07-09",
	}

	err = repo.CreateActor(actor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Before Update: %+v\n", actor)

	err = repo.DeleteActors(5)
	// actor.Name = "Yeab Hanks"

	// err = repo.Update(actor)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// updatedActor, err := repo.FindById(actor.Id)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("After Update: %+v\n", actor)
}
