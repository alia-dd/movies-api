package handler

import (
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"

	"net/http"
	"strconv"
)

/*
POST - createGenre
GET - getAllGenres
GET - getGenresByID
PUT - updateGenre
Delete -deleteGenreByName/ByID
*/
type GenreHandler struct {
	genreService *service.GenreService
}

func NewGenreHandler(genreService *service.GenreService) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var genre models.Genre
	/*
	   r.Body-- is raw json from request
	   NewDecoder parses json
	   Decode(&genre) converts Json into go struct and populates genre

	*/
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.genreService.CreateGenre(&genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.genreService.GetAllGenres()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) GetGenreByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid genre id", http.StatusBadRequest)
		return
	}

	genre, err := h.genreService.GetGenreByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var genre models.Genre

	err = json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.genreService.UpdateGenre(id, genre.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre Updated"})

}

func (h *GenreHandler) DeleteGenreByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.genreService.DeleteGenreByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre Deleted"})
}

func (h *GenreHandler) DeleteGenreByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	err := h.genreService.DeleteGenreByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Genre Deleted"})
}
