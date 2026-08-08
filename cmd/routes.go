package main

import (
	"movies-api/internal/database"
	"movies-api/internal/handler"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"
)

type dbData struct {
	dbName          string
	movieEntityName string
}

func route() (http.Handler, error) {
	mux := http.NewServeMux()

	db, dbErr := database.InitializeDB("movie.db")
	if dbErr != nil {
		return nil, dbErr
	}

	movieTableErr := database.InitializeMovieTable(db)
	if movieTableErr != nil {
		return nil, movieTableErr
	}

	movieRepo := repository.NewdbConnection(db)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handler.NewHandler(movieService)

	// these endpoints call the movies entity handler
	mux.HandleFunc("POST /api/movies", handler.PostMovie)
	mux.HandleFunc("GET /api/movies", movieHandler.GetMovies)

	mux.HandleFunc("GET /api/movies/{id}", movieHandler.GetMoviesById)

	mux.HandleFunc("PATCH /api/movies/{id}", handler.PatchMovie)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.DeleteMovie)

	mux.HandleFunc("GET /api/movies/{movieId}/actors", movieHandler.GetActorsForMovie) // get all actor in selected movie

	return mux, nil
}

// Set up the following endpoints for each entity:

// POST /api/{entity}: Create a new entity
// GET /api/{entity}: Retrieve all entities
// GET /api/{entity}/{id}: Retrieve a specific entity by ID
// PATCH /api/{entity}/{id}: Partially update an existing entity
// DELETE /api/{entity}/{id}: Delete an entity
// Additionally, implement filtering endpoints for the following:

// GET /api/movies?genre={genreId}: Retrieve movies filtered by genre
// GET /api/movies?year={releaseYear}: Retrieve movies filtered by release year
// GET /api/movies?actor={actorId}: Retrieve movies that the actor with the given id has starred in
// GET /api/movies/{movieId}/actors: Retrieve all actors starring in a movie
// GET /api/actors?name={name}:
