package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/back/events"
	"github.com/back/kafkax"
	"github.com/back/models"
	"github.com/back/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RoomHandler struct {
	repo   repository.RoomRepository
	events *kafkax.Publisher
}

// Función constructora
func NewRoomHandler(repo repository.RoomRepository, events *kafkax.Publisher) *RoomHandler {
	return &RoomHandler{repo: repo, events: events}
}

func (h *RoomHandler) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	inmuebleID := r.URL.Query().Get("inmuebleId") // Da parámetros del query '/back?inmuebleId=10', y se extrae roomId

	var rooms []models.Room
	var err error

	if inmuebleID != "" {
		rooms, err = h.repo.GetByInmuebleID(inmuebleID)
	} else {
		rooms, err = h.repo.GetAll()
	}

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
	// Verificar que el room exista (el evento necesita sus IDs)
	room, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // Error 404
		w.Write([]byte("Room not found"))
		return
	}

	if err := h.repo.Delete(params["id"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Error 500
		w.Write([]byte(err.Error()))
		return
	}

	// Avisa a servicios para que borre los servicios de la habitación
	evt := events.RoomDeleted{
		EventID:    uuid.NewString(),
		RoomID:     room.ID,
		InmuebleID: room.InmuebleID,
		DeletedAt:  time.Now().UTC(),
	}
	key := strconv.FormatUint(uint64(room.InmuebleID), 10)
	if err := h.events.Publish(r.Context(), key, events.TypeRoomDeleted, evt); err != nil {
		log.Printf("no se pudo publicar %s de la habitación %d: %v", events.TypeRoomDeleted, room.ID, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Room deleted"))
}
