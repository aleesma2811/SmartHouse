package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/servicios/models"
	"github.com/servicios/repository"
)

type PlugHandler struct {
	repo repository.PlugRepository
}

func NewPlugHandler(repo repository.PlugRepository) *PlugHandler {
	return &PlugHandler{repo: repo}
}

func (h *PlugHandler) GetPlugsHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId") // Da parámetros del query '/plugs?roomId=10', y se extrae roomId

	var plugs []models.Plug
	var err error

	if roomID != "" {
		plugs, err = h.repo.GetByRoomID(roomID)
	} else {
		plugs, err = h.repo.GetAll()
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&plugs)
}

func (h *PlugHandler) GetPlugHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	plug, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plug not found"))
		return
	}
	json.NewEncoder(w).Encode(plug)
}

func (h *PlugHandler) PostPlugsHandler(w http.ResponseWriter, r *http.Request) {
	var plug models.Plug

	if err := json.NewDecoder(r.Body).Decode(&plug); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.repo.Create(&plug); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(plug)
}

func (h *PlugHandler) UpdatePlugHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	plug, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plug not found"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(plug); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.repo.Update(plug); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(plug)
}

func (h *PlugHandler) DeletePlugsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if _, err := h.repo.GetByID(params["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.repo.Delete(params["id"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Plug deleted"))
}
