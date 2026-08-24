package handler

import (
	"database/sql"
	"encoding/json"
	"movies-api/internal/errors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ActorHandler struct {
	service *service.ActorService
}

func NewActorHandler(service *service.ActorService) *ActorHandler {
	return &ActorHandler{
		service: service,
	}
}

func (h *ActorHandler) CreateActor(w http.ResponseWriter, r *http.Request) {

	var actor models.Actor

	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	err = h.service.CreateActor(&actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid actor id", http.StatusBadRequest)
		return
	}
	var actor models.Actor

	err = json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	actor.Id = id
	err = h.service.UpdateActor(&actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" {
		pageStr = "1"
	}
	if limitStr == "" {
		limitStr = "10"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if page < 1 || limit < 1 {
		http.Error(w, "page and limit must be greater than 0", http.StatusBadRequest)
		return
	}
	actors, total, err := h.service.GetAllActors(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data": actors,
		"pagination": map[string]int{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

func (h *ActorHandler) GetActorsById(w http.ResponseWriter, r *http.Request) {

	Id, errId := strconv.Atoi(r.PathValue("id"))
	if errId != nil {
		http.Error(w, "Invalid actor id", http.StatusBadRequest)
		return
	}
	actors, err := h.service.FindById(Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) GetActorByName(w http.ResponseWriter, r *http.Request) {

	name := strings.TrimSpace(r.URL.Query().Get("name"))

	actor, err := h.service.FindByName(name)

	if err == sql.ErrNoRows {
		http.Error(w, errors.ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) DeleteActorsById(w http.ResponseWriter, r *http.Request) {

	Id, errId := strconv.Atoi(r.PathValue("id"))
	if errId != nil {
		http.Error(w, "Invalid actor id", http.StatusBadRequest)
		return
	}
	force := r.URL.Query().Get("force") == "true"
	err := h.service.DeleteActorsById(Id, force)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *ActorHandler) DeleteActorsByName(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")
	force := r.URL.Query().Get("force") == "true"
	err := h.service.DeleteActorsByName(name, force)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
