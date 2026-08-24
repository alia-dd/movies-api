package main

import (
	"database/sql"
	"flag"
	"fmt"
	"movies-api/internal/database"
	"movies-api/internal/errors"
	"movies-api/internal/handler"
	"movies-api/internal/repository"
	"movies-api/internal/service"
	"net/http"
	"strings"
)

type dbData struct {
	dbName          string
	movieEntityName string
}

func route() (http.Handler, error) {
	seed := flag.Bool("s", false, "Use flag -s to Seed the Database.")
	generate := flag.String("g", "", "Use flag -g to Generate the ApiKey.")

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

	if *generate != "" {
		key, err := database.GenerateApiKey(db, generate)
		if err != nil {
			return nil, err
		}
		fmt.Println(key)
	}

	//actor routes
	repo := repository.NewActorRepository(db)
	actorService := service.NewActorService(repo)
	actorHandler := handler.NewActorHandler(actorService)

	mux.HandleFunc("POST /api/actors", Middleware(db, actorHandler.CreateActor))
	mux.HandleFunc("PATCH /api/actors/{id}", Middleware(db, actorHandler.UpdateActor))
	mux.HandleFunc("GET /api/actors", Middleware(db, actorHandler.GetAllActors))
	mux.HandleFunc("GET /api/actors/{id}", Middleware(db, actorHandler.GetActorsById))
	mux.HandleFunc("GET /api/actors/name", Middleware(db, actorHandler.GetActorByName))
	mux.HandleFunc("DELETE /api/actors/{id}", Middleware(db, actorHandler.DeleteActorsById))
	mux.HandleFunc("DELETE /api/actors/name/{name}", Middleware(db, actorHandler.DeleteActorsByName))

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handler.NewMovieHandler(movieService)

	// Movie routes
	mux.HandleFunc("POST /api/movies", Middleware(db, movieHandler.CreateMovie))
	mux.HandleFunc("GET /api/movies", Middleware(db, movieHandler.GetAllMovies))
	mux.HandleFunc("GET /api/movies/{id}", Middleware(db, movieHandler.GetMoviesById))
	mux.HandleFunc("GET /api/movies/search", Middleware(db, movieHandler.GetMoviesByTitle))
	mux.HandleFunc("PATCH /api/movies/{id}", Middleware(db, movieHandler.UpdateMovie))
	mux.HandleFunc("DELETE /api/movies/{id}", Middleware(db, movieHandler.DeleteMovie))
	mux.HandleFunc("GET /api/movies/{movieId}/actors", Middleware(db, movieHandler.GetActorsForMovie)) // get all actor in selected movie
	mux.HandleFunc("GET /api/movies/{movieId}/genres", Middleware(db, movieHandler.GetGenresForMovie)) // get all actor in selected movie

	// Genre routes
	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	mux.HandleFunc("POST /api/genres", Middleware(db, genreHandler.CreateGenre))
	mux.HandleFunc("GET /api/genres", Middleware(db, genreHandler.GetAllGenres))
	mux.HandleFunc("GET /api/genres/search", Middleware(db, genreHandler.SearchGenreByName))
	mux.HandleFunc("GET /api/genres/{id}", Middleware(db, genreHandler.GetGenreByID))
	mux.HandleFunc("PATCH /api/genres/{id}", Middleware(db, genreHandler.UpdateGenre))
	mux.HandleFunc("DELETE /api/genres/{id}", Middleware(db, genreHandler.DeleteGenreByID))
	mux.HandleFunc("DELETE /api/genres/name/{name}", Middleware(db, genreHandler.DeleteGenreByName))
	return mux, nil
}

func Middleware(db *sql.DB, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, errors.ErrServerErr.Error(), http.StatusInternalServerError)
			}
		}()
		key := strings.TrimSpace(r.URL.Query().Get("apikey"))
		if found := database.GetApiKey(db, key); !found {
			http.Error(w, errors.ErrWrongApiKey.Error(), http.StatusBadRequest)
			return
		}
		handler(w, r)
	}
}
