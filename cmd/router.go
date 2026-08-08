package main

import (
	"log"
	"moviesApi/internal/database"
	"moviesApi/internal/repository"
	"moviesApi/internal/service"
	"moviesApi/internal/handler"
	"net/http"
)

func route() (http.Handler, error) {
	mux := http.NewServeMux()

	db, err := database.OpenDB("movies.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	err = database.NewTable(db)

	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewActorRepository(db)

	actorService := service.NewActorService(repo)

	actorHandler := handler.NewActorHandler(actorService)

	mux.HandleFunc("POST /actors", actorHandler.CreateActor)
	mux.HandleFunc("GET /actors", actorHandler.GetAllActors)
	mux.HandleFunc("GET/actors/{id}", actorHandler.GetActorsById)
	mux.HandleFunc("GET/actors/name/{name}", actorHandler.GetActorsById)

	return mux, nil
}
