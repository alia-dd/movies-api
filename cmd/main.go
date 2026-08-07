package main

import (
	"fmt"
	"log"
	"moviesApi/internal/database"
	"moviesApi/internal/models"
	"moviesApi/internal/repository"
	"moviesApi/internal/service"
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
	actorService := service.NewActorService(repo)
	actor := &models.Actor{
		Name:      "",
		BirthDate: "1978-07-09",
	}
	err = actorService.CreateActor(actor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("data: %+v\n", actor)

	// err = repo.DeleteActors(5)
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
