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
	defer db.Close()

	err = database.NewTable(db)

	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewActorRepository(db)

	actorService := service.NewActorService(repo)

	actorHandler := handler.NewActorHandler(actorService)

	mux.HandleFunc("POST /api/actors", actorHandler.CreateActor)
	mux.HandleFunc("PUT /api/actors/{id}", actorHandler.UpdateActor)
	mux.HandleFunc("GET /api/actors", actorHandler.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetActorsById)
	mux.HandleFunc("GET /api/actors/name/{name}", actorHandler.GetActorByName)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.DeleteActorsById)
	mux.HandleFunc("DELETE /api/actors/name/{name}", actorHandler.DeleteActorsByName)

	return mux, nil
}
