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
	
	err = database.NewTable(db)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	//actor routes
	actorRepo := repository.NewActorRepository(db)
	actorService := service.NewActorService(actorRepo)
	actorHandler := handler.NewActorHandler(actorService)

	mux.HandleFunc("POST /api/actors", actorHandler.CreateActor)
	mux.HandleFunc("PUT /api/actors/{id}", actorHandler.UpdateActor)
	mux.HandleFunc("GET /api/actors", actorHandler.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetActorsById)
	mux.HandleFunc("GET /api/actors/name/{name}", actorHandler.GetActorByName)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.DeleteActorsById)
	mux.HandleFunc("DELETE /api/actors/name/{name}", actorHandler.DeleteActorsByName)
	
	// Genre routes
	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	mux.HandleFunc("POST /api/genres", genreHandler.CreateGenre)
	mux.HandleFunc("GET /api/genres", genreHandler.GetAllGenres)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GetGenreByID)
	mux.HandleFunc("PUT /api/genres/{id}", genreHandler.UpdateGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.DeleteGenreByID)
	mux.HandleFunc("DELETE /api/genres/name/{name}", genreHandler.DeleteGenreByName)

	return mux, nil
}
