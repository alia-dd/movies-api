package handler

import (
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
)

type GenreHandler struct {
	genreService *service.GenreService
}

func NewGenreHandler(genreService *service.GenreService) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	var genre models.Genre

	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	err = h.genreService.CreateGenre(cx, &genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "invalid page parameter", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "invalid limit parameter", http.StatusBadRequest)
		return
	}

	if page <= 0 || limit <= 0 {
		http.Error(w, "page and limit must be greater than 0", http.StatusBadRequest)
		return
	}

	genres, total, err := h.genreService.GetAllGenres(cx, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": genres,
		"pagination": map[string]int{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

func (h *GenreHandler) GetGenreByID(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid genre ID", http.StatusBadRequest)
		return
	}

	genre, err := h.genreService.GetGenreByID(cx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) SearchGenreByName(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	searchTerm := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	if searchTerm == "" {
		http.Error(w, "search query cannot be empty", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "invalid page parameter", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "invalid limit parameter", http.StatusBadRequest)
		return
	}

	if page <= 0 || limit <= 0 {
		http.Error(w, "page and limit must be greater than 0", http.StatusBadRequest)
		return
	}

	genres, total, err := h.genreService.SearchGenreByName(cx, searchTerm, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": genres,
		"pagination": map[string]int{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

func (h *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid genre ID", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err = h.genreService.UpdateGenre(cx, id, updateData.Name)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Genre not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	genre, _ := h.genreService.GetGenreByID(cx, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) DeleteGenreByID(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid genre ID", http.StatusBadRequest)
		return
	}

	forceStr := r.URL.Query().Get("force")
	force := forceStr == "true"

	err = h.genreService.DeleteGenreByID(cx, id, force)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}
		if !force {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) DeleteGenreByName(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	forceStr := r.URL.Query().Get("force")
	force := forceStr == "true"

	err := h.genreService.DeleteGenreByName(cx, name, force)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}
		if !force {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
