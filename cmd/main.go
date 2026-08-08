package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := ":8000"
	fmt.Println("server running on localhost", port)
	srv, err := route()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(port, srv))

	// err = database.NewTable(db)
	// if err != nil {
	// 	log.Fatal(err)

	// }

	// fmt.Println("table created")

	// repo := repository.NewActorRepository(db)
	// actorService := service.NewActorService(repo)
	// actor := &models.Actor{
	// 	Name:      "test2",
	// 	BirthDate: "1978-07-09",
	// }
	// err = actorService.CreateActor(actor)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("data: %+v\n", actor)

	// err = repo.DeleteActors(5)
	// actor.Name = "Yeab Hanks"

	// err = repo.Update(actor)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = repo.DeleteActorsByName("test1")

	// updatedActor, err := repo.FindByName("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("After Update: %+v\n", actor)
}
