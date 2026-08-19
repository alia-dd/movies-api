package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllMovies(w http.ResponseWriter, r *http.Request) {

	var f models.Filter
	f.Genre = strings.TrimSpace(r.URL.Query().Get("genre"))
	f.Actor = strings.TrimSpace(r.URL.Query().Get("actor"))
	f.Year = strings.TrimSpace(r.URL.Query().Get("year"))
	f.Page = strings.TrimSpace(r.URL.Query().Get("page"))
	f.Size = strings.TrimSpace(r.URL.Query().Get("size"))

	payload, err := h.service.GetMovie(f)
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}

func (h *Handler) GetMoviesById(w http.ResponseWriter, r *http.Request) {
	var payload *models.MoviesDisplay
	var byTitlePayload []models.MoviesDisplay
	var err error
	id, idErr := strconv.Atoi(r.PathValue("id"))
	if idErr != nil {
		byTitlePayload, err = h.service.SearchMovieByTitle(r.PathValue("id"))
	} else {
		payload, err = h.service.GetMovieById(id)
	}

	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var jsonData []byte
	var jsonErr error
	if byTitlePayload != nil {
		jsonData, jsonErr = json.MarshalIndent(byTitlePayload, "", "  ")
	} else {
		jsonData, jsonErr = json.MarshalIndent(payload, "", "  ")
	}
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}

// add new movie to the db
func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movies
	jsonErr := json.NewDecoder(r.Body).Decode(&movie)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	lastInsertedMovie, err := h.service.CreateMovie(movie)
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(lastInsertedMovie, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}

// update movie with given id with the given data
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte(id))
}

// delete the movie with given id
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	val := r.URL.Query().Get("force")
	if val == "" {
		val = "false"
	}
	sentencedId := r.PathValue("id")
	if deleteErr := h.service.DeleteMovie(sentencedId, val); deleteErr != nil {
		if deleteErr == sql.ErrNoRows {
			http.Error(w, deleteErr.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, deleteErr.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal([]string{"message: ", "movie deleted"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}

func (h *Handler) GetGenresForMovie(w http.ResponseWriter, r *http.Request) {

	var payload []string
	var err error
	id, idErr := strconv.Atoi(r.PathValue("movieId"))
	if idErr != nil {
		movie, movieErr := h.service.GetMovieByTitle(r.PathValue("movieId"))
		if movieErr != nil {
			messge := "Invalid Movie ID"
			if movieErr == sql.ErrNoRows {
				messge = "Movie Not Found"
			}
			jsonData, err := json.Marshal([]string{"message: ", messge})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonData)
			return
		}
		id = movie.Id
	}
	payload, err = h.service.GetGenresForMovie(id)

	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, jsonErr := json.MarshalIndent(payload, "", "  ")

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}
func (h *Handler) GetActorsForMovie(w http.ResponseWriter, r *http.Request) {

	var payload []string
	var err error
	id, idErr := strconv.Atoi(r.PathValue("movieId"))
	if idErr != nil {
		movie, movieErr := h.service.GetMovieByTitle(r.PathValue("movieId"))
		if movieErr != nil {
			messge := "Invalid Movie ID"
			if movieErr == sql.ErrNoRows {
				messge = "Movie Not Found"
			}
			jsonData, err := json.Marshal([]string{"message: ", messge})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonData)
			return
		}
		id = movie.Id
	}
	payload, err = h.service.GetActorsForMovie(id)

	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, jsonErr := json.MarshalIndent(payload, "", "  ")

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}
