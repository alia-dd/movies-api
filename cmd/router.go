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

	db, dbErr := database.InitializeDB("internal/database/data/moviesApi.db")
	if dbErr != nil {
		return nil, dbErr
	}

	movieTableErr := database.InitializeMovieTable(db)
	if movieTableErr != nil {
		return nil, movieTableErr
	}

	//actor routes
	repo := repository.NewActorRepository(db)
	actorService := service.NewActorService(repo)
	actorHandler := handler.NewActorHandler(actorService)

	mux.HandleFunc("POST /api/actors", actorHandler.CreateActor)
	mux.HandleFunc("PATCH /api/actors/{id}", actorHandler.UpdateActor)
	mux.HandleFunc("GET /api/actors", actorHandler.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetActorsById)
	mux.HandleFunc("GET /api/actors/name/{name}", actorHandler.GetActorByName)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.DeleteActorsById)
	mux.HandleFunc("DELETE /api/actors/name/{name}", actorHandler.DeleteActorsByName)

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handler.NewMovieHandler(movieService)

	// Movie routes
	mux.HandleFunc("POST /api/movies", movieHandler.CreateMovie)
	mux.HandleFunc("GET /api/movies", movieHandler.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}", movieHandler.GetMoviesById)
	mux.HandleFunc("PATCH /api/movies/{id}", movieHandler.UpdateMovie)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.DeleteMovie)
	mux.HandleFunc("GET /api/movies/{movieId}/actors", movieHandler.GetActorsForMovie) // get all actor in selected movie
	mux.HandleFunc("GET /api/movies/{movieId}/genres", movieHandler.GetGenresForMovie) // get all actor in selected movie

	// Genre routes
	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	mux.HandleFunc("POST /api/genres", genreHandler.CreateGenre)
	mux.HandleFunc("GET /api/genres", genreHandler.GetAllGenres)
	mux.HandleFunc("GET /api/genres/search", genreHandler.SearchGenreByName)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GetGenreByID)
	mux.HandleFunc("PATCH /api/genres/{id}", genreHandler.UpdateGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.DeleteGenreByID)
	mux.HandleFunc("DELETE /api/genres/name/{name}", genreHandler.DeleteGenreByName)
	return mux, nil
}
