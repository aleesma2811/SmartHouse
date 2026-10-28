package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/inmuebles/events"
	"github.com/inmuebles/kafkax"
	"github.com/inmuebles/models"
	"github.com/inmuebles/repository"
)

type InmuebleHandler struct {
	repo   repository.InmuebleRepository
	events *kafkax.Publisher
}

func NewInmuebleHandler(repo repository.InmuebleRepository, events *kafkax.Publisher) *InmuebleHandler {
	return &InmuebleHandler{repo: repo, events: events}
}

func tipoValidado(tipo string) bool {
	switch tipo {
	case models.TipoCasa, models.TipoDepartamento:
		return true
	}
	return false
}

func (h *InmuebleHandler) GetInmueblesHandler(w http.ResponseWriter, r *http.Request) {
	inmuebles, err := h.repo.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&inmuebles)
}

func (h *InmuebleHandler) GetInmuebleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inmueble, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Inmueble not found"))
		return
	}
	json.NewEncoder(w).Encode(inmueble)
}

func (h *InmuebleHandler) PostInmueblesHandler(w http.ResponseWriter, r *http.Request) {
	var inmueble models.Inmueble

	if err := json.NewDecoder(r.Body).Decode(&inmueble); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !tipoValidado(inmueble.Tipo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tipo inválido: debe ser 'casa' o 'departamento'"))
		return
	}

	if err := h.repo.Create(&inmueble); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&inmueble)
}

func (h *InmuebleHandler) UpdateInmueble(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inmueble, err := h.repo.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Servicio not found"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(inmueble); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !tipoValidado(inmueble.Tipo) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Tipo inválido: debe ser 'casa' o 'departamento'"))
		return
	}

	if err := h.repo.Update(inmueble); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(inmueble)
}

func (h *InmuebleHandler) DeleteInmuebleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Evento necesita el ID
	inmueble, err := h.repo.GetByID(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.repo.Delete(params["id"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Crea el evento y genera el ID del evento
	// Kafka avisa que inmueble fue borrado, back reacciona eliminando todas las habitaciones
	evt := events.InmuebleDeleted{
		EventID:    uuid.NewString(),
		InmuebleID: inmueble.ID,
		DeletedAt:  time.Now().UTC(),
	}

	if err := h.events.Publish(r.Context(), params["id"], events.TypeInmuebleDeleted, evt); err != nil {
		log.Printf("no se pudo publicar %s del inmueble %d: %v", events.TypeInmuebleDeleted, inmueble.ID, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sevicio deleted"))
}
