package handler

import (
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
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

	actors, err := h.service.GetAllActors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
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

	name := r.PathValue("name")

	actor, err := h.service.FindByName(name)

	if err != nil {
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

	err := h.service.DeleteActorsById(Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *ActorHandler) DeleteActorsByName(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")

	err := h.service.DeleteActorsByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
