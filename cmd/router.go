package main

import (
	"flag"
	"fmt"
	"movies-api/internal/database"
	"movies-api/internal/errors"
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
	seed := flag.Bool("s", false, "Use flag -s to Seed the Database.")
	flag.Parse()
	mux := http.NewServeMux()

	db, dbErr := database.InitializeDB("internal/database/data/moviesApi.db")
	if dbErr != nil {
		return nil, dbErr
	}

	movieTableErr := database.InitializeMovieTable(db)
	if movieTableErr != nil {
		return nil, movieTableErr
	}

	if *seed {
		if err := database.SeedData(db); err != nil {
			return nil, fmt.Errorf("failed to seed movie data: %w", err)
		}
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
	mux.HandleFunc("POST /api/movies", Middleware(movieHandler.CreateMovie))
	mux.HandleFunc("GET /api/movies", Middleware(movieHandler.GetAllMovies))
	mux.HandleFunc("GET /api/movies/{id}", Middleware(movieHandler.GetMoviesById))
	mux.HandleFunc("PATCH /api/movies/{id}", Middleware(movieHandler.UpdateMovie))
	mux.HandleFunc("DELETE /api/movies/{id}", Middleware(movieHandler.DeleteMovie))
	mux.HandleFunc("GET /api/movies/{movieId}/actors", Middleware(movieHandler.GetActorsForMovie)) // get all actor in selected movie
	mux.HandleFunc("GET /api/movies/{movieId}/genres", Middleware(movieHandler.GetGenresForMovie)) // get all actor in selected movie

	// Genre routes
	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	mux.HandleFunc("POST /api/genres", Middleware(genreHandler.CreateGenre))
	mux.HandleFunc("GET /api/genres", Middleware(genreHandler.GetAllGenres))
	mux.HandleFunc("GET /api/genres/search", Middleware(genreHandler.SearchGenreByName))
	mux.HandleFunc("GET /api/genres/{id}", Middleware(genreHandler.GetGenreByID))
	mux.HandleFunc("PATCH /api/genres/{id}", Middleware(genreHandler.UpdateGenre))
	mux.HandleFunc("DELETE /api/genres/{id}", Middleware(genreHandler.DeleteGenreByID))
	mux.HandleFunc("DELETE /api/genres/name/{name}", Middleware(genreHandler.DeleteGenreByName))
	return mux, nil
}

func Middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, errors.ErrServerErr.Error(), http.StatusInternalServerError)
			}
		}()
		handler(w, r)
	}
}
