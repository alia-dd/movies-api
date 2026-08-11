package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *Handler {
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
	var payload *models.Movies
	var byTitlePayload []models.Movies
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
func (h *Handler) PostMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movies
	jsonErr := json.NewDecoder(r.Body).Decode(&movie)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateMovie(movie)
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(movie, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
}

// update movie with given id with the given data
func PatchMovie(w http.ResponseWriter, r *http.Request) {
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
