package handler

import (
	"database/sql"
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.MovieService
}

func NewHandler(service *service.MovieService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {

	var f models.Filter
	f.Genre = r.URL.Query().Get("genre")
	f.Actor = r.URL.Query().Get("actor")
	f.Year = r.URL.Query().Get("year")
	f.Page = r.URL.Query().Get("page")
	f.Size = r.URL.Query().Get("size")

	payload, err := h.service.GetMovie(f)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (h *Handler) GetMoviesById(w http.ResponseWriter, r *http.Request) {
	var payload *models.Movies
	var err error
	id, idErr := strconv.Atoi(r.PathValue("id"))
	if idErr != nil {
		payload, err = h.service.GetMovieByTitle(r.PathValue("id"))
	} else {
		payload, err = h.service.GetMovieById(id)
	}

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func GetMoviesByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	w.Write([]byte(title))
}

// add new movie to the db
func PostMovie(w http.ResponseWriter, r *http.Request) {

}

// update movie with given id with the given data
func PatchMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte(id))
}

// delete the movie with given id
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	val := r.URL.Query().Get("force")

	id, idErr := strconv.Atoi(r.PathValue("id"))
	if idErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if deleteErr := h.service.DeleteMovie(id, val); deleteErr != nil {
		if deleteErr == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetActorsForMovie(w http.ResponseWriter, r *http.Request) {
	// id, idErr := strconv.Atoi(r.PathValue("id"))
	// if idErr != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// payload, err := h.service.(id)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// w.WriteHeader(http.StatusAccepted)
	// jsonData, err := json.Marshal(payload)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(jsonData)
}
