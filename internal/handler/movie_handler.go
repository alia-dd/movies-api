package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movies-api/internal/errors"
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
	cx := r.Context()
	var f models.Filter
	f.Genre = strings.TrimSpace(r.URL.Query().Get("genre"))
	f.Actor = strings.TrimSpace(r.URL.Query().Get("actor"))
	f.Year = strings.TrimSpace(r.URL.Query().Get("year"))
	f.Page = strings.TrimSpace(r.URL.Query().Get("page"))
	f.Size = strings.TrimSpace(r.URL.Query().Get("size"))

	payload, total, err := h.service.GetMovie(cx, f)
	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(f.Page)
	size, err := strconv.Atoi(f.Size)
	jsonData, err := json.MarshalIndent(map[string]interface{}{
		"pagination": map[string]int{
			"page":        page,
			"limit":       size,
			"total":       total,
			"total_pages": (total + size - 1) / size,
		},
		"data": payload}, "", "  ",
	)
	if err != nil {
		http.Error(w, errors.ErrorMarshel.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetMoviesById(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	var payload *models.MoviesDisplay
	var byTitlePayload []models.MoviesDisplay
	var err error

	id, idErr := strconv.Atoi(r.PathValue("id"))
	if idErr != nil {
		byTitlePayload, err = h.service.SearchMovieByTitle(cx, r.PathValue("id"))
	} else {
		payload, err = h.service.GetMovieById(cx, id)
	}

	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var jsonData []byte
	var jsonErr error
	if byTitlePayload != nil {
		jsonData, jsonErr = json.MarshalIndent(map[string]interface{}{"data": byTitlePayload}, "", "  ")
	} else {
		jsonData, jsonErr = json.MarshalIndent(map[string]interface{}{"data": payload}, "", "  ")
	}
	if jsonErr != nil {
		http.Error(w, errors.ErrorMarshel.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// add new movie to the db
func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	var movie models.Movies
	jsonErr := json.NewDecoder(r.Body).Decode(&movie)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		http.Error(w, errors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	lastInsertedMovie, err := h.service.CreateMovie(cx, movie)
	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(map[string]interface{}{"data": lastInsertedMovie}, "", "  ")
	if err != nil {
		http.Error(w, errors.ErrorMarshel.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// update movie with given id with the given data
func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	id, idErr := strconv.Atoi(r.PathValue("id"))
	if idErr != nil {
		http.Error(w, errors.ErrInvalidId.Error(), http.StatusBadRequest)
		return
	}

	var movie models.MovieUpdate
	jsonErr := json.NewDecoder(r.Body).Decode(&movie)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		http.Error(w, errors.ErrInvalidInput.Error(), http.StatusBadRequest)
		return
	}

	lastUpdatedMovie, err := h.service.UpdateMovie(cx, id, movie)
	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(map[string]interface{}{"data": lastUpdatedMovie}, "", "  ")
	if err != nil {
		http.Error(w, errors.ErrorMarshel.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// delete the movie with given id
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	val := r.URL.Query().Get("force")
	if val == "" {
		val = "false"
	}
	sentencedId := r.PathValue("id")
	if deleteErr := h.service.DeleteMovie(cx, sentencedId, val); deleteErr != nil {
		if deleteErr == sql.ErrNoRows {
			http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, deleteErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetGenresForMovie(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	var payload []string
	var err error

	id, idErr := strconv.Atoi(r.PathValue("movieId"))
	if idErr != nil {
		movie, movieErr := h.service.GetMovieByTitle(cx, r.PathValue("movieId"))
		if movieErr != nil {
			if movieErr == sql.ErrNoRows {
				http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}
		id = movie.Id
	}

	payload, err = h.service.GetGenresForMovie(cx, id)
	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, jsonErr := json.MarshalIndent(map[string]interface{}{"data": payload}, "", "  ")
	if jsonErr != nil {
		http.Error(w, errors.ErrorMarshel.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetActorsForMovie(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	var payload []string
	var err error

	id, idErr := strconv.Atoi(r.PathValue("movieId"))
	if idErr != nil {
		movie, movieErr := h.service.GetMovieByTitle(cx, r.PathValue("movieId"))
		if movieErr != nil {

			if movieErr == sql.ErrNoRows {
				http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}
		id = movie.Id
	}

	payload, err = h.service.GetActorsForMovie(cx, id)
	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, jsonErr := json.MarshalIndent(map[string]interface{}{"data": payload}, "", "  ")
	if jsonErr != nil {
		http.Error(w, errors.ErrorMarshel.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
