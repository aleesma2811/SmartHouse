package routes

import (
	"encoding/json"
	"net/http"

	"github.com/back/models"
	"github.com/back/repository"
	"github.com/gorilla/mux"
)

type RoomHandler struct {
	repo repository.RoomRepository
}

// Función constructora
func NewRoomHandler(repo repository.RoomRepository) *RoomHandler {
	return &RoomHandler{repo: repo}
}

func (h *RoomHandler) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.repo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&rooms)
}

func (h *RoomHandler) GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	room, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Room not found"))
		return
	}
	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) PostRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	// Si el JSON viene mal formado
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Si algo más falla
	if err := h.repo.Create(&room); err != nil {
		w.WriteHeader(http.StatusBadRequest) // Error 400
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&room)
}

func (h *RoomHandler) DeleteRoomsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verificar que el room exista
	if _, err := h.repo.GetByID(params["id"]); err != nil {
		w.WriteHeader(http.StatusNotFound) // Error 404
		w.Write([]byte("Room not found"))
		return
	}

	if err := h.repo.Delete(params["id"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Error 500
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Room deleted"))
}
